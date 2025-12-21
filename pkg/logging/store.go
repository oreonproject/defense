// oreon/defense · watchthelight <wtl>

package logging

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "modernc.org/sqlite"
)

// LogEntry represents a log entry in the database.
type LogEntry struct {
	ID        int64
	Timestamp time.Time
	Level     string
	Component string
	Message   string
	Metadata  map[string]interface{}
}

// LogStore manages log storage in SQLite.
type LogStore struct {
	db *sql.DB
}

// NewLogStore creates a new log store.
// Use ":memory:" for path to create an in-memory database (useful for tests).
func NewLogStore(path string) (*LogStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &LogStore{db: db}, nil
}

// Close closes the database connection.
func (s *LogStore) Close() error {
	return s.db.Close()
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		level TEXT NOT NULL,
		component TEXT,
		message TEXT NOT NULL,
		metadata TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
	`
	_, err := db.Exec(schema)
	return err
}

// Insert adds a log entry to the database.
func (s *LogStore) Insert(entry LogEntry) error {
	var metadataJSON []byte
	var err error

	if len(entry.Metadata) > 0 {
		metadataJSON, err = json.Marshal(entry.Metadata)
		if err != nil {
			return err
		}
	}

	_, err = s.db.Exec(
		`INSERT INTO logs (timestamp, level, component, message, metadata) VALUES (?, ?, ?, ?, ?)`,
		entry.Timestamp, entry.Level, entry.Component, entry.Message, metadataJSON,
	)
	return err
}

// QueryOptions specifies filters for querying logs.
type QueryOptions struct {
	Level     string    // filter by level (empty = all)
	Component string    // filter by component (empty = all)
	After     time.Time // only logs after this time (zero = no filter)
	Before    time.Time // only logs before this time (zero = no filter)
	Limit     int       // max results (0 = no limit)
}

// Query retrieves log entries matching the options.
// Results ordered by timestamp desc (newest first).
func (s *LogStore) Query(opts QueryOptions) ([]LogEntry, error) {
	query := `SELECT id, timestamp, level, component, message, metadata FROM logs WHERE 1=1`
	args := []interface{}{}

	if opts.Level != "" {
		query += ` AND level = ?`
		args = append(args, opts.Level)
	}
	if opts.Component != "" {
		query += ` AND component = ?`
		args = append(args, opts.Component)
	}
	if !opts.After.IsZero() {
		query += ` AND timestamp > ?`
		args = append(args, opts.After)
	}
	if !opts.Before.IsZero() {
		query += ` AND timestamp < ?`
		args = append(args, opts.Before)
	}

	query += ` ORDER BY timestamp DESC`

	if opts.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, opts.Limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LogEntry
	for rows.Next() {
		var entry LogEntry
		var metadataJSON []byte

		err := rows.Scan(&entry.ID, &entry.Timestamp, &entry.Level, &entry.Component, &entry.Message, &metadataJSON)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &entry.Metadata)
		}

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// Prune deletes log entries older than the given duration.
// Returns the number of entries deleted.
func (s *LogStore) Prune(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)

	result, err := s.db.Exec(`DELETE FROM logs WHERE timestamp < ?`, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
