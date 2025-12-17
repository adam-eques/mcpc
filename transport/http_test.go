package transport

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPRequestResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "initialize") {
			t.Errorf("unexpected body: %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{}}`))
	}))
	defer srv.Close()

	tr := NewHTTP(srv.URL)
	ctx := context.Background()
	if err := tr.Send(ctx, []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`)); err != nil {
		t.Fatal(err)
	}
	reply, err := tr.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(reply), `"result"`) {
		t.Fatalf("reply=%s", reply)
	}
}

func TestHTTPNotificationNoReply(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	tr := NewHTTP(srv.URL)
	if err := tr.Send(context.Background(), []byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)); err != nil {
		t.Fatal(err)
	}
	if len(tr.replies) != 0 {
		t.Fatalf("notification should not enqueue a reply, queued %d", len(tr.replies))
	}
}

func TestHTTPErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	if err := NewHTTP(srv.URL).Send(context.Background(), []byte(`{}`)); err == nil {
		t.Fatal("expected error for 500 response")
	}
}
