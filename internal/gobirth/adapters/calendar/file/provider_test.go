package file

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestProvider_EventsForDate_FiltersByDate(t *testing.T) {
	tmp, err := os.CreateTemp("", "gobirth-events-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer os.Remove(tmp.Name())

	json := `[
  {"id":"1","title":"Pepe","description":"phone: +34600111222","start_date":"2026-01-16"},
  {"id":"2","title":"Ana","description":"phone: +34600333444","start_date":"2026-01-17"}
]`
	if _, err := tmp.WriteString(json); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	_ = tmp.Close()

	p := Provider{Path: tmp.Name()}

	date := time.Date(2026, 1, 16, 12, 0, 0, 0, time.UTC)
	events, err := p.EventsForDate(context.Background(), date)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Title != "Pepe" {
		t.Fatalf("expected Pepe, got %q", events[0].Title)
	}
}
