package server

import "testing"

func TestNewHertz_CreatesServer(t *testing.T) {
	hs := NewHertz(`function handle(req,res){ res.setStatus(200); res.end('ok'); }`)
	if hs == nil {
		t.Fatalf("expected server instance")
	}
}
