package client

import (
	"encoding/json"
	"testing"

	"github.com/adam-eques/mcpc/mcp"
)

func TestProgressTrackerRoutes(t *testing.T) {
	tr := NewProgressTracker()
	got := make(chan mcp.ProgressParams, 1)
	tr.Track("job-1", func(p mcp.ProgressParams) { got <- p })

	params, _ := json.Marshal(mcp.ProgressParams{ProgressToken: "job-1", Progress: 0.5, Total: 1})
	tr.Handle(mcp.NotificationProgress, params)

	select {
	case p := <-got:
		if p.Progress != 0.5 {
			t.Fatalf("progress=%v", p.Progress)
		}
	default:
		t.Fatal("callback not invoked")
	}
}

func TestProgressTrackerIgnoresOthers(t *testing.T) {
	tr := NewProgressTracker()
	// Should not panic or block on a non-progress notification.
	tr.Handle(mcp.NotificationMessage, json.RawMessage(`{}`))
	tr.Forget("nope")
}
