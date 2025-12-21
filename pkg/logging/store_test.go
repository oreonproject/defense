// oreon/defense · watchthelight <wtl>

package logging

import (
	"testing"
	"time"
)

func TestNewLogStore(t *testing.T) {
	store, err := NewLogStore(":memory:")
	if err != nil {
		t.Fatalf("NewLogStore: %v", err)
	}
	defer store.Close()
}

func TestInsertAndQuery(t *testing.T) {
	store, err := NewLogStore(":memory:")
	if err != nil {
		t.Fatalf("NewLogStore: %v", err)
	}
	defer store.Close()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Component: "test",
		Message:   "hello world",
		Metadata:  map[string]interface{}{"key": "value"},
	}

	if err := store.Insert(entry); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	results, err := store.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Message != "hello world" {
		t.Errorf("message = %q, want %q", results[0].Message, "hello world")
	}
	if results[0].Metadata["key"] != "value" {
		t.Errorf("metadata[key] = %v, want %q", results[0].Metadata["key"], "value")
	}
}

func TestQueryFilterLevel(t *testing.T) {
	store, err := NewLogStore(":memory:")
	if err != nil {
		t.Fatalf("NewLogStore: %v", err)
	}
	defer store.Close()

	store.Insert(LogEntry{Timestamp: time.Now(), Level: "info", Message: "info msg"})
	store.Insert(LogEntry{Timestamp: time.Now(), Level: "error", Message: "error msg"})
	store.Insert(LogEntry{Timestamp: time.Now(), Level: "info", Message: "info msg 2"})

	results, err := store.Query(QueryOptions{Level: "error"})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Message != "error msg" {
		t.Errorf("message = %q, want %q", results[0].Message, "error msg")
	}
}

func TestQueryFilterComponent(t *testing.T) {
	store, err := NewLogStore(":memory:")
	if err != nil {
		t.Fatalf("NewLogStore: %v", err)
	}
	defer store.Close()

	store.Insert(LogEntry{Timestamp: time.Now(), Level: "info", Component: "scanner", Message: "scan msg"})
	store.Insert(LogEntry{Timestamp: time.Now(), Level: "info", Component: "daemon", Message: "daemon msg"})

	results, err := store.Query(QueryOptions{Component: "scanner"})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Component != "scanner" {
		t.Errorf("component = %q, want %q", results[0].Component, "scanner")
	}
}

func TestQueryLimit(t *testing.T) {
	store, err := NewLogStore(":memory:")
	if err != nil {
		t.Fatalf("NewLogStore: %v", err)
	}
	defer store.Close()

	for i := 0; i < 10; i++ {
		store.Insert(LogEntry{Timestamp: time.Now(), Level: "info", Message: "msg"})
	}

	results, err := store.Query(QueryOptions{Limit: 3})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestPrune(t *testing.T) {
	store, err := NewLogStore(":memory:")
	if err != nil {
		t.Fatalf("NewLogStore: %v", err)
	}
	defer store.Close()

	old := time.Now().Add(-48 * time.Hour)
	recent := time.Now()

	store.Insert(LogEntry{Timestamp: old, Level: "info", Message: "old"})
	store.Insert(LogEntry{Timestamp: old, Level: "info", Message: "old 2"})
	store.Insert(LogEntry{Timestamp: recent, Level: "info", Message: "recent"})

	deleted, err := store.Prune(24 * time.Hour)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if deleted != 2 {
		t.Errorf("deleted = %d, want 2", deleted)
	}

	results, _ := store.Query(QueryOptions{})
	if len(results) != 1 {
		t.Errorf("expected 1 remaining, got %d", len(results))
	}
}
