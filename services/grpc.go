package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// GrpcService handles gRPC requests.
type GrpcService struct{}

type GrpcServiceInfo struct {
	Name    string       `json:"name"`
	Methods []GrpcMethod `json:"methods"`
}

type GrpcMethod struct {
	Name           string `json:"name"`
	FullName       string `json:"fullName"`
	InputType      string `json:"inputType"`
	OutputType     string `json:"outputType"`
	IsClientStream bool   `json:"isClientStream"`
	IsServerStream bool   `json:"isServerStream"`
	InputTemplate  string `json:"inputTemplate"`
}

type GrpcResponse struct {
	Data     string            `json:"data"`
	Metadata map[string]string `json:"metadata"`
	Duration int64             `json:"duration"`
	Error    string            `json:"error,omitempty"`
}

// ParseProtoFile parses .proto file(s) and returns services/methods with input templates.
func (g *GrpcService) ParseProtoFile(filePath string, importPaths []string) ([]GrpcServiceInfo, error) {
	if len(importPaths) == 0 {
		// Default: use the directory of the proto file
		parts := strings.Split(filePath, "/")
		if len(parts) > 1 {
			importPaths = []string{strings.Join(parts[:len(parts)-1], "/")}
			filePath = parts[len(parts)-1]
		} else {
			importPaths = []string{"."}
		}
	}

	parser := protoparse.Parser{
		ImportPaths: importPaths,
	}

	fds, err := parser.ParseFiles(filePath)
	if err != nil {
		return nil, fmt.Errorf("parse proto: %w", err)
	}

	var services []GrpcServiceInfo
	for _, fd := range fds {
		for _, svc := range fd.GetServices() {
			info := GrpcServiceInfo{Name: svc.GetFullyQualifiedName()}
			for _, method := range svc.GetMethods() {
				template := buildMessageTemplate(method.GetInputType())
				info.Methods = append(info.Methods, GrpcMethod{
					Name:           method.GetName(),
					FullName:       fmt.Sprintf("/%s/%s", svc.GetFullyQualifiedName(), method.GetName()),
					InputType:      method.GetInputType().GetFullyQualifiedName(),
					OutputType:     method.GetOutputType().GetFullyQualifiedName(),
					IsClientStream: method.IsClientStreaming(),
					IsServerStream: method.IsServerStreaming(),
					InputTemplate:  template,
				})
			}
			services = append(services, info)
		}
	}
	return services, nil
}

// SendUnary sends a unary gRPC request using a proto file for type information.
func (g *GrpcService) SendUnary(address, protoFile string, importPaths []string, fullMethod, messageJSON string, metadataPairs []Pair) GrpcResponse {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Add metadata
	md := metadata.New(nil)
	for _, p := range metadataPairs {
		if p.Enabled && p.Name != "" {
			md.Append(p.Name, p.Value)
		}
	}
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Connect
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return GrpcResponse{Error: fmt.Sprintf("connect: %v", err), Duration: time.Since(start).Milliseconds()}
	}
	defer conn.Close()

	// Parse proto
	if len(importPaths) == 0 {
		importPaths = []string{"."}
	}
	parser := protoparse.Parser{ImportPaths: importPaths}
	fds, err := parser.ParseFiles(protoFile)
	if err != nil {
		return GrpcResponse{Error: fmt.Sprintf("parse proto: %v", err), Duration: time.Since(start).Milliseconds()}
	}

	// Find method descriptor
	var methodDesc *desc.MethodDescriptor
	for _, fd := range fds {
		for _, svc := range fd.GetServices() {
			for _, m := range svc.GetMethods() {
				mFullName := fmt.Sprintf("/%s/%s", svc.GetFullyQualifiedName(), m.GetName())
				if mFullName == fullMethod {
					methodDesc = m
					break
				}
			}
			if methodDesc != nil {
				break
			}
		}
		if methodDesc != nil {
			break
		}
	}

	if methodDesc == nil {
		return GrpcResponse{Error: "method not found: " + fullMethod, Duration: time.Since(start).Milliseconds()}
	}

	// Build request message
	reqMsg := dynamic.NewMessage(methodDesc.GetInputType())
	if err := reqMsg.UnmarshalJSON([]byte(messageJSON)); err != nil {
		return GrpcResponse{Error: fmt.Sprintf("unmarshal input: %v", err), Duration: time.Since(start).Milliseconds()}
	}

	// Invoke
	stub := grpcdynamic.NewStub(conn)
	var headerMD, trailerMD metadata.MD
	respMsg, err := stub.InvokeRpc(ctx, methodDesc, reqMsg,
		grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return GrpcResponse{Error: fmt.Sprintf("rpc error: %v", err), Duration: duration}
	}

	// Marshal response
	dynResp, ok := respMsg.(*dynamic.Message)
	if !ok {
		return GrpcResponse{Error: "unexpected response type", Duration: duration}
	}
	respJSON, _ := dynResp.MarshalJSON()

	// Collect response metadata
	respMeta := make(map[string]string)
	for k, v := range headerMD {
		respMeta[k] = strings.Join(v, ", ")
	}
	for k, v := range trailerMD {
		respMeta["trailer-"+k] = strings.Join(v, ", ")
	}

	return GrpcResponse{
		Data:     string(respJSON),
		Metadata: respMeta,
		Duration: duration,
	}
}

// buildMessageTemplate creates a JSON template from a message descriptor.
func buildMessageTemplate(msg *desc.MessageDescriptor) string {
	result := make(map[string]interface{})
	for _, field := range msg.GetFields() {
		name := field.GetJSONName()
		if name == "" {
			name = field.GetName()
		}
		switch field.GetType().String() {
		case "TYPE_STRING":
			result[name] = ""
		case "TYPE_INT32", "TYPE_INT64", "TYPE_SINT32", "TYPE_SINT64",
			"TYPE_UINT32", "TYPE_UINT64", "TYPE_FIXED32", "TYPE_FIXED64",
			"TYPE_SFIXED32", "TYPE_SFIXED64":
			result[name] = 0
		case "TYPE_FLOAT", "TYPE_DOUBLE":
			result[name] = 0.0
		case "TYPE_BOOL":
			result[name] = false
		case "TYPE_BYTES":
			result[name] = ""
		case "TYPE_ENUM":
			result[name] = 0
		case "TYPE_MESSAGE":
			result[name] = json.RawMessage("{}")
		default:
			result[name] = nil
		}
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	return string(b)
}
