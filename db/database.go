package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize opens or creates the SQLite database in the user's config directory.
func Initialize() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dbDir := filepath.Join(configDir, "gapiro")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("create db dir: %w", err)
	}

	dbPath := filepath.Join(dbDir, "gapiro.db")
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL")
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite single-writer
	db.SetMaxIdleConns(1)
	DB = db

	return migrate()
}

func migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT 'Default',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS folders (
			id TEXT PRIMARY KEY,
			workspace_id TEXT NOT NULL,
			parent_id TEXT,
			name TEXT NOT NULL,
			sort_order INTEGER DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS requests (
			id TEXT PRIMARY KEY,
			workspace_id TEXT NOT NULL,
			folder_id TEXT,
			name TEXT NOT NULL DEFAULT 'Untitled',
			method TEXT NOT NULL DEFAULT 'GET',
			url TEXT NOT NULL DEFAULT '',
			headers TEXT NOT NULL DEFAULT '[]',
			url_parameters TEXT NOT NULL DEFAULT '[]',
			body_type TEXT NOT NULL DEFAULT 'none',
			body TEXT NOT NULL DEFAULT '',
			form_data TEXT NOT NULL DEFAULT '[]',
			auth TEXT NOT NULL DEFAULT '{"type":"none"}',
			settings TEXT NOT NULL DEFAULT '{}',
			description TEXT NOT NULL DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS responses (
			id TEXT PRIMARY KEY,
			request_id TEXT NOT NULL,
			status INTEGER NOT NULL DEFAULT 0,
			status_text TEXT NOT NULL DEFAULT '',
			headers TEXT NOT NULL DEFAULT '{}',
			body TEXT NOT NULL DEFAULT '',
			size INTEGER NOT NULL DEFAULT 0,
			duration INTEGER NOT NULL DEFAULT 0,
			dns_time INTEGER NOT NULL DEFAULT 0,
			connect_time INTEGER NOT NULL DEFAULT 0,
			tls_time INTEGER NOT NULL DEFAULT 0,
			ttfb_time INTEGER NOT NULL DEFAULT 0,
			error TEXT NOT NULL DEFAULT '',
			protocol TEXT NOT NULL DEFAULT '',
			remote_addr TEXT NOT NULL DEFAULT '',
			content_type TEXT NOT NULL DEFAULT '',
			redirect_count INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS environments (
			id TEXT PRIMARY KEY,
			workspace_id TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT 'Default',
			variables TEXT NOT NULL DEFAULT '[]',
			is_active INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS grpc_requests (
			id TEXT PRIMARY KEY,
			workspace_id TEXT NOT NULL,
			folder_id TEXT,
			name TEXT NOT NULL DEFAULT 'Untitled',
			url TEXT NOT NULL DEFAULT '',
			proto_files TEXT NOT NULL DEFAULT '[]',
			service TEXT NOT NULL DEFAULT '',
			method TEXT NOT NULL DEFAULT '',
			message TEXT NOT NULL DEFAULT '',
			metadata TEXT NOT NULL DEFAULT '[]',
			sort_order INTEGER DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		)`,
		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_requests_workspace ON requests(workspace_id)`,
		`CREATE INDEX IF NOT EXISTS idx_responses_request ON responses(request_id)`,
		`CREATE INDEX IF NOT EXISTS idx_folders_workspace ON folders(workspace_id)`,
	}

	for _, m := range migrations {
		if _, err := DB.Exec(m); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}

	// Ensure a default workspace exists
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM workspaces").Scan(&count)
	if count == 0 {
		now := time.Now().UnixMilli()
		DB.Exec("INSERT INTO workspaces (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
			"ws_default", "Default Workspace", now, now)
	}

	return nil
}

// Close shuts down the database connection.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// Helper: marshal to JSON string
func ToJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Helper: unmarshal from JSON string
func FromJSON(s string, v interface{}) {
	json.Unmarshal([]byte(s), v)
}
