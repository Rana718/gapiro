package services

import (
	"changeme/db"
	"encoding/json"
	"fmt"
	"time"
)

// CollectionService manages workspaces, folders, requests, and responses (CRUD + persistence).
type CollectionService struct{}

// ─── Types ────────────────────────────────────────────────────────────────────

type Pair struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type AuthConfig struct {
	Type   string          `json:"type"`
	Basic  json.RawMessage `json:"basic,omitempty"`
	Bearer json.RawMessage `json:"bearer,omitempty"`
	APIKey json.RawMessage `json:"apiKey,omitempty"`
}

type RequestSettings struct {
	Timeout         int  `json:"timeout"`
	FollowRedirects bool `json:"followRedirects"`
	VerifySSL       bool `json:"verifySSL"`
	MaxRedirects    int  `json:"maxRedirects"`
}

type SavedRequest struct {
	ID            string          `json:"id"`
	WorkspaceID   string          `json:"workspaceId"`
	FolderID      string          `json:"folderId,omitempty"`
	Name          string          `json:"name"`
	Method        string          `json:"method"`
	URL           string          `json:"url"`
	Headers       []Pair          `json:"headers"`
	URLParameters []Pair          `json:"urlParameters"`
	BodyType      string          `json:"bodyType"`
	Body          string          `json:"body"`
	FormData      []Pair          `json:"formData"`
	Auth          AuthConfig      `json:"auth"`
	Settings      RequestSettings `json:"settings"`
	Description   string          `json:"description"`
	SortOrder     int             `json:"sortOrder"`
	CreatedAt     int64           `json:"createdAt"`
	UpdatedAt     int64           `json:"updatedAt"`
}

type SavedResponse struct {
	ID            string `json:"id"`
	RequestID     string `json:"requestId"`
	Status        int    `json:"status"`
	StatusText    string `json:"statusText"`
	Headers       string `json:"headers"`
	Body          string `json:"body"`
	Size          int64  `json:"size"`
	Duration      int64  `json:"duration"`
	DNSTime       int64  `json:"dnsTime"`
	ConnectTime   int64  `json:"connectTime"`
	TLSTime       int64  `json:"tlsTime"`
	TTFBTime      int64  `json:"ttfbTime"`
	Error         string `json:"error"`
	Protocol      string `json:"protocol"`
	RemoteAddr    string `json:"remoteAddr"`
	ContentType   string `json:"contentType"`
	RedirectCount int    `json:"redirectCount"`
	CreatedAt     int64  `json:"createdAt"`
}

type Folder struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	ParentID    string `json:"parentId,omitempty"`
	Name        string `json:"name"`
	SortOrder   int    `json:"sortOrder"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type EnvVariable struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type Environment struct {
	ID          string        `json:"id"`
	WorkspaceID string        `json:"workspaceId"`
	Name        string        `json:"name"`
	Variables   []EnvVariable `json:"variables"`
	IsActive    bool          `json:"isActive"`
	CreatedAt   int64         `json:"createdAt"`
	UpdatedAt   int64         `json:"updatedAt"`
}

// ─── Workspace CRUD ───────────────────────────────────────────────────────────

func (c *CollectionService) ListWorkspaces() ([]Workspace, error) {
	rows, err := db.DB.Query("SELECT id, name, created_at, updated_at FROM workspaces ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Workspace
	for rows.Next() {
		var w Workspace
		rows.Scan(&w.ID, &w.Name, &w.CreatedAt, &w.UpdatedAt)
		result = append(result, w)
	}
	return result, nil
}

func (c *CollectionService) CreateWorkspace(name string) (Workspace, error) {
	now := time.Now().UnixMilli()
	id := generateID("ws")
	_, err := db.DB.Exec(
		"INSERT INTO workspaces (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		id, name, now, now,
	)
	if err != nil {
		return Workspace{}, err
	}
	return Workspace{ID: id, Name: name, CreatedAt: now, UpdatedAt: now}, nil
}

// ─── Folder CRUD ──────────────────────────────────────────────────────────────

func (c *CollectionService) ListFolders(workspaceID string) ([]Folder, error) {
	rows, err := db.DB.Query(
		"SELECT id, workspace_id, parent_id, name, sort_order, created_at, updated_at FROM folders WHERE workspace_id = ? ORDER BY sort_order, name",
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Folder
	for rows.Next() {
		var f Folder
		var parentID *string
		rows.Scan(&f.ID, &f.WorkspaceID, &parentID, &f.Name, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt)
		if parentID != nil {
			f.ParentID = *parentID
		}
		result = append(result, f)
	}
	return result, nil
}

func (c *CollectionService) CreateFolder(workspaceID, name, parentID string) (Folder, error) {
	now := time.Now().UnixMilli()
	id := generateID("fl")
	var pID *string
	if parentID != "" {
		pID = &parentID
	}
	_, err := db.DB.Exec(
		"INSERT INTO folders (id, workspace_id, parent_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		id, workspaceID, pID, name, now, now,
	)
	if err != nil {
		return Folder{}, err
	}
	return Folder{ID: id, WorkspaceID: workspaceID, ParentID: parentID, Name: name, CreatedAt: now, UpdatedAt: now}, nil
}

func (c *CollectionService) DeleteFolder(id string) error {
	_, err := db.DB.Exec("DELETE FROM folders WHERE id = ?", id)
	return err
}

// ─── Request CRUD ─────────────────────────────────────────────────────────────

func (c *CollectionService) ListRequests(workspaceID string) ([]SavedRequest, error) {
	rows, err := db.DB.Query(
		`SELECT id, workspace_id, folder_id, name, method, url, headers, url_parameters,
		 body_type, body, form_data, auth, settings, description, sort_order, created_at, updated_at
		 FROM requests WHERE workspace_id = ? ORDER BY sort_order, updated_at DESC`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SavedRequest
	for rows.Next() {
		var r SavedRequest
		var folderID *string
		var headersJSON, paramsJSON, formDataJSON, authJSON, settingsJSON string

		rows.Scan(&r.ID, &r.WorkspaceID, &folderID, &r.Name, &r.Method, &r.URL,
			&headersJSON, &paramsJSON, &r.BodyType, &r.Body, &formDataJSON,
			&authJSON, &settingsJSON, &r.Description, &r.SortOrder, &r.CreatedAt, &r.UpdatedAt)

		if folderID != nil {
			r.FolderID = *folderID
		}
		json.Unmarshal([]byte(headersJSON), &r.Headers)
		json.Unmarshal([]byte(paramsJSON), &r.URLParameters)
		json.Unmarshal([]byte(formDataJSON), &r.FormData)
		json.Unmarshal([]byte(authJSON), &r.Auth)
		json.Unmarshal([]byte(settingsJSON), &r.Settings)

		if r.Headers == nil {
			r.Headers = []Pair{}
		}
		if r.URLParameters == nil {
			r.URLParameters = []Pair{}
		}
		if r.FormData == nil {
			r.FormData = []Pair{}
		}

		result = append(result, r)
	}
	if result == nil {
		result = []SavedRequest{}
	}
	return result, nil
}

func (c *CollectionService) SaveRequest(req SavedRequest) (SavedRequest, error) {
	now := time.Now().UnixMilli()

	if req.ID == "" {
		req.ID = generateID("rq")
		req.CreatedAt = now
	}
	req.UpdatedAt = now

	if req.WorkspaceID == "" {
		req.WorkspaceID = "ws_default"
	}
	if req.Headers == nil {
		req.Headers = []Pair{}
	}
	if req.URLParameters == nil {
		req.URLParameters = []Pair{}
	}
	if req.FormData == nil {
		req.FormData = []Pair{}
	}

	headersJSON, _ := json.Marshal(req.Headers)
	paramsJSON, _ := json.Marshal(req.URLParameters)
	formDataJSON, _ := json.Marshal(req.FormData)
	authJSON, _ := json.Marshal(req.Auth)
	settingsJSON, _ := json.Marshal(req.Settings)

	var folderID *string
	if req.FolderID != "" {
		folderID = &req.FolderID
	}

	_, err := db.DB.Exec(`
		INSERT OR REPLACE INTO requests
		(id, workspace_id, folder_id, name, method, url, headers, url_parameters,
		 body_type, body, form_data, auth, settings, description, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ID, req.WorkspaceID, folderID, req.Name, req.Method, req.URL,
		string(headersJSON), string(paramsJSON), req.BodyType, req.Body,
		string(formDataJSON), string(authJSON), string(settingsJSON),
		req.Description, req.SortOrder, req.CreatedAt, req.UpdatedAt,
	)
	if err != nil {
		return SavedRequest{}, err
	}
	return req, nil
}

func (c *CollectionService) DeleteRequest(id string) error {
	_, err := db.DB.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}

func (c *CollectionService) DuplicateRequest(id string) (SavedRequest, error) {
	reqs, _ := c.ListRequests("ws_default")
	for _, r := range reqs {
		if r.ID == id {
			r.ID = generateID("rq")
			r.Name = r.Name + " (copy)"
			r.CreatedAt = time.Now().UnixMilli()
			r.UpdatedAt = r.CreatedAt
			return c.SaveRequest(r)
		}
	}
	return SavedRequest{}, fmt.Errorf("request not found: %s", id)
}

// ─── Response History ─────────────────────────────────────────────────────────

func (c *CollectionService) SaveResponse(requestID string, resp SavedResponse) (SavedResponse, error) {
	resp.ID = generateID("rs")
	resp.RequestID = requestID
	resp.CreatedAt = time.Now().UnixMilli()

	_, err := db.DB.Exec(`
		INSERT INTO responses
		(id, request_id, status, status_text, headers, body, size, duration,
		 dns_time, connect_time, tls_time, ttfb_time, error, protocol, remote_addr,
		 content_type, redirect_count, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		resp.ID, resp.RequestID, resp.Status, resp.StatusText, resp.Headers,
		resp.Body, resp.Size, resp.Duration, resp.DNSTime, resp.ConnectTime,
		resp.TLSTime, resp.TTFBTime, resp.Error, resp.Protocol, resp.RemoteAddr,
		resp.ContentType, resp.RedirectCount, resp.CreatedAt,
	)
	return resp, err
}

func (c *CollectionService) ListResponses(requestID string, limit int) ([]SavedResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := db.DB.Query(
		`SELECT id, request_id, status, status_text, headers, body, size, duration,
		 dns_time, connect_time, tls_time, ttfb_time, error, protocol, remote_addr,
		 content_type, redirect_count, created_at
		 FROM responses WHERE request_id = ? ORDER BY created_at DESC LIMIT ?`,
		requestID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SavedResponse
	for rows.Next() {
		var r SavedResponse
		rows.Scan(&r.ID, &r.RequestID, &r.Status, &r.StatusText, &r.Headers,
			&r.Body, &r.Size, &r.Duration, &r.DNSTime, &r.ConnectTime,
			&r.TLSTime, &r.TTFBTime, &r.Error, &r.Protocol, &r.RemoteAddr,
			&r.ContentType, &r.RedirectCount, &r.CreatedAt)
		result = append(result, r)
	}
	if result == nil {
		result = []SavedResponse{}
	}
	return result, nil
}

func (c *CollectionService) ClearResponseHistory(requestID string) error {
	_, err := db.DB.Exec("DELETE FROM responses WHERE request_id = ?", requestID)
	return err
}

// ─── Environment CRUD ─────────────────────────────────────────────────────────

func (c *CollectionService) ListEnvironments(workspaceID string) ([]Environment, error) {
	rows, err := db.DB.Query(
		"SELECT id, workspace_id, name, variables, is_active, created_at, updated_at FROM environments WHERE workspace_id = ? ORDER BY name",
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Environment
	for rows.Next() {
		var e Environment
		var varsJSON string
		var isActive int
		rows.Scan(&e.ID, &e.WorkspaceID, &e.Name, &varsJSON, &isActive, &e.CreatedAt, &e.UpdatedAt)
		json.Unmarshal([]byte(varsJSON), &e.Variables)
		e.IsActive = isActive == 1
		if e.Variables == nil {
			e.Variables = []EnvVariable{}
		}
		result = append(result, e)
	}
	if result == nil {
		result = []Environment{}
	}
	return result, nil
}

func (c *CollectionService) SaveEnvironment(env Environment) (Environment, error) {
	now := time.Now().UnixMilli()
	if env.ID == "" {
		env.ID = generateID("ev")
		env.CreatedAt = now
	}
	env.UpdatedAt = now
	if env.WorkspaceID == "" {
		env.WorkspaceID = "ws_default"
	}
	if env.Variables == nil {
		env.Variables = []EnvVariable{}
	}

	varsJSON, _ := json.Marshal(env.Variables)
	isActive := 0
	if env.IsActive {
		isActive = 1
	}

	_, err := db.DB.Exec(`
		INSERT OR REPLACE INTO environments (id, workspace_id, name, variables, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		env.ID, env.WorkspaceID, env.Name, string(varsJSON), isActive, env.CreatedAt, env.UpdatedAt,
	)
	return env, err
}

func (c *CollectionService) DeleteEnvironment(id string) error {
	_, err := db.DB.Exec("DELETE FROM environments WHERE id = ?", id)
	return err
}

func (c *CollectionService) SetActiveEnvironment(workspaceID, envID string) error {
	// Deactivate all in workspace
	db.DB.Exec("UPDATE environments SET is_active = 0 WHERE workspace_id = ?", workspaceID)
	// Activate the chosen one
	if envID != "" {
		db.DB.Exec("UPDATE environments SET is_active = 1 WHERE id = ?", envID)
	}
	return nil
}
