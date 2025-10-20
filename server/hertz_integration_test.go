package server

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func startServer(t *testing.T, script string) *http.Client {
	hs := NewHertz(script)
	go hs.Spin()
	// give the server a brief moment to start
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = hs.Shutdown(ctx)
	})
	time.Sleep(150 * time.Millisecond)
	return &http.Client{Timeout: 2 * time.Second}
}

func TestHertz_Proxy_SuccessAndError(t *testing.T) {
	// Single server instance with conditional behavior to avoid port conflicts
	script := `
	  function handle(req, res){
	    if (req.path === '/err') { throw new Error('bad'); }
	    res.setHeader('X-Ok','1');
	    res.setStatus(201);
	    res.end('hi ' + req.path);
	  }
	`
	client := startServer(t, script)
	// success
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/abc?x=1", nil)
	req.Header.Set("X-Req", "v")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Ok") != "1" {
		t.Fatalf("missing X-Ok header")
	}
	if string(b) != "hi /abc" {
		t.Fatalf("unexpected body: %q", string(b))
	}
	// error
	resp2, err := client.Get("http://127.0.0.1:8888/err")
	if err != nil {
		t.Fatalf("request error2: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 500 {
		t.Fatalf("expected 500, got %d", resp2.StatusCode)
	}
}
