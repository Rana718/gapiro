package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptrace"
	"os"
	"strings"
	"sync"
	"time"
)

// HttpService handles all HTTP request execution from the frontend.
type HttpService struct {
	mu      sync.Mutex
	cancels map[string]context.CancelFunc
}

// KeyValue represents a key-value pair for headers, params, form data.
type KeyValue struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Enabled   bool   `json:"enabled"`
	ValueType string `json:"valueType,omitempty"`
	FileName  string `json:"fileName,omitempty"`
}

// RequestPayload is the request configuration sent from the frontend.
type RequestPayload struct {
	RequestID       string     `json:"requestId"`
	Method          string     `json:"method"`
	URL             string     `json:"url"`
	Headers         []KeyValue `json:"headers"`
	QueryParams     []KeyValue `json:"queryParams"`
	BodyType        string     `json:"bodyType"` // none, json, text, form-urlencoded, form-data, binary
	Body            string     `json:"body"`
	FormData        []KeyValue `json:"formData"`
	Timeout         int        `json:"timeout"` // seconds, 0 = no timeout
	FollowRedirects bool       `json:"followRedirects"`
	VerifySSL       bool       `json:"verifySSL"`
	MaxRedirects    int        `json:"maxRedirects"`
}

// ResponsePayload is the response data sent back to the frontend.
type ResponsePayload struct {
	Status        int               `json:"status"`
	StatusText    string            `json:"statusText"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	Size          int64             `json:"size"`
	Duration      float64           `json:"duration"` // milliseconds (fractional for sub-ms)
	DNSTime       float64           `json:"dnsTime"`
	ConnectTime   float64           `json:"connectTime"`
	TLSTime       float64           `json:"tlsTime"`
	TTFBTime      float64           `json:"ttfbTime"` // time to first byte
	Error         string            `json:"error,omitempty"`
	Protocol      string            `json:"protocol"`
	RemoteAddr    string            `json:"remoteAddr"`
	ContentType   string            `json:"contentType"`
	RedirectCount int               `json:"redirectCount"`
}

// SendRequest executes an HTTP request and returns the response with timing details.
// The context is provided by Wails and cancelled when the frontend calls .cancel() on the promise.
func (h *HttpService) SendRequest(ctx context.Context, payload RequestPayload) ResponsePayload {
	emptyResponse := func() ResponsePayload {
		return ResponsePayload{Headers: map[string]string{}, Body: ""}
	}
	requestCtx, cancel := context.WithCancel(ctx)
	if payload.RequestID != "" {
		h.mu.Lock()
		if h.cancels == nil {
			h.cancels = make(map[string]context.CancelFunc)
		}
		h.cancels[payload.RequestID] = cancel
		h.mu.Unlock()
		defer func() {
			h.mu.Lock()
			delete(h.cancels, payload.RequestID)
			h.mu.Unlock()
		}()
	}
	defer cancel()
	// Build URL with query params
	url := payload.URL
	if !strings.Contains(url, "://") {
		url = "http://" + url
	}

	enabledParams := []KeyValue{}
	for _, p := range payload.QueryParams {
		if p.Enabled && p.Key != "" {
			enabledParams = append(enabledParams, p)
		}
	}
	if len(enabledParams) > 0 {
		separator := "?"
		if strings.Contains(url, "?") {
			separator = "&"
		}
		params := []string{}
		for _, p := range enabledParams {
			params = append(params, fmt.Sprintf("%s=%s", p.Key, p.Value))
		}
		url += separator + strings.Join(params, "&")
	}

	// Build request body
	var body io.Reader
	var contentType string

	switch payload.BodyType {
	case "json":
		body = strings.NewReader(payload.Body)
		contentType = "application/json"
	case "text":
		body = strings.NewReader(payload.Body)
		contentType = "text/plain"
	case "form-urlencoded":
		formValues := []string{}
		for _, kv := range payload.FormData {
			if kv.Enabled && kv.Key != "" {
				formValues = append(formValues, fmt.Sprintf("%s=%s", kv.Key, kv.Value))
			}
		}
		body = strings.NewReader(strings.Join(formValues, "&"))
		contentType = "application/x-www-form-urlencoded"
	case "form-data":
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		for _, kv := range payload.FormData {
			if kv.Enabled && kv.Key != "" {
				if kv.ValueType == "file" {
					file, err := os.Open(kv.Value)
					if err != nil {
						continue
					}
					name := kv.FileName
					if name == "" {
						parts := strings.Split(strings.ReplaceAll(kv.Value, "\\", "/"), "/")
						name = parts[len(parts)-1]
					}
					part, err := writer.CreateFormFile(kv.Key, name)
					if err == nil {
						_, _ = io.Copy(part, file)
					}
					_ = file.Close()
				} else {
					_ = writer.WriteField(kv.Key, kv.Value)
				}
			}
		}
		writer.Close()
		body = &buf
		contentType = writer.FormDataContentType()
	case "binary":
		if payload.Body != "" {
			data, err := os.ReadFile(payload.Body)
			if err == nil {
				body = bytes.NewReader(data)
				contentType = "application/octet-stream"
			}
		}
	}

	// Create request
	req, err := http.NewRequest(payload.Method, url, body)
	if err != nil {
		result := emptyResponse()
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		return result
	}

	// Set content type if we have a body
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Set headers
	for _, h := range payload.Headers {
		if h.Enabled && h.Key != "" {
			req.Header.Set(h.Key, h.Value)
		}
	}

	// Set default User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Gapiro/1.0")
	}

	// HTTP trace for timing
	var dnsStart, dnsEnd, connectStart, connectEnd, tlsStart, tlsEnd, gotFirstByte time.Time
	var remoteAddr string
	redirectCount := 0

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			dnsEnd = time.Now()
		},
		ConnectStart: func(_, _ string) {
			connectStart = time.Now()
		},
		ConnectDone: func(_, addr string, _ error) {
			connectEnd = time.Now()
			remoteAddr = addr
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			tlsEnd = time.Now()
		},
		GotFirstResponseByte: func() {
			gotFirstByte = time.Now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(requestCtx, trace))

	// Configure client
	timeout := time.Duration(30) * time.Second
	if payload.Timeout > 0 {
		timeout = time.Duration(payload.Timeout) * time.Second
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: timeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !payload.VerifySSL,
		},
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	if !payload.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			redirectCount = len(via)
			maxRedirects := payload.MaxRedirects
			if maxRedirects <= 0 {
				maxRedirects = 10
			}
			if len(via) >= maxRedirects {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		}
	}

	// Send request
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		if requestCtx.Err() != nil {
			result := emptyResponse()
			result.Error = "Request cancelled"
			result.Duration = duration.Seconds() * 1000
			return result
		}
		result := emptyResponse()
		result.Error = fmt.Sprintf("Request failed: %v", err)
		result.Duration = duration.Seconds() * 1000
		return result
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponsePayload{
			Status:     resp.StatusCode,
			StatusText: resp.Status,
			Error:      fmt.Sprintf("Failed to read response body: %v", err),
			Duration:   duration.Seconds() * 1000,
		}
	}

	// Build response headers map
	headers := make(map[string]string)
	for key, values := range resp.Header {
		headers[key] = strings.Join(values, ", ")
	}

	// Calculate timing
	var dnsTime, connectTime, tlsTime, ttfbTime float64
	if !dnsStart.IsZero() && !dnsEnd.IsZero() {
		dnsTime = dnsEnd.Sub(dnsStart).Seconds() * 1000
	}
	if !connectStart.IsZero() && !connectEnd.IsZero() {
		connectTime = connectEnd.Sub(connectStart).Seconds() * 1000
	}
	if !tlsStart.IsZero() && !tlsEnd.IsZero() {
		tlsTime = tlsEnd.Sub(tlsStart).Seconds() * 1000
	}
	if !gotFirstByte.IsZero() {
		ttfbTime = gotFirstByte.Sub(start).Seconds() * 1000
	}

	// Determine protocol version
	protocol := fmt.Sprintf("HTTP/%d.%d", resp.ProtoMajor, resp.ProtoMinor)

	return ResponsePayload{
		Status:        resp.StatusCode,
		StatusText:    resp.Status,
		Headers:       headers,
		Body:          string(respBody),
		Size:          int64(len(respBody)),
		Duration:      duration.Seconds() * 1000,
		DNSTime:       dnsTime,
		ConnectTime:   connectTime,
		TLSTime:       tlsTime,
		TTFBTime:      ttfbTime,
		Protocol:      protocol,
		RemoteAddr:    remoteAddr,
		ContentType:   resp.Header.Get("Content-Type"),
		RedirectCount: redirectCount,
	}
}

// CancelRequest explicitly aborts a native request. This works independently
// of the frontend bridge promise and is reliable during DNS/TLS/body reads.
func (h *HttpService) CancelRequest(requestID string) bool {
	h.mu.Lock()
	cancel := h.cancels[requestID]
	h.mu.Unlock()
	if cancel == nil {
		return false
	}
	cancel()
	return true
}

// FormatJSON takes a JSON string and returns it pretty-printed.
func (h *HttpService) FormatJSON(input string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return input
	}
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return input
	}
	return string(formatted)
}
