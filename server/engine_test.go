package server

import (
	"strings"
	"testing"
)

const sampleScript = `
function handle(req, res) {
  res.setHeader('Content-Type', 'application/json');
  res.setStatus(201);
  // echo some fields from request into body
  const out = {
    method: req.method,
    path: req.path,
    url: req.url,
    query: req.query,
    headers: req.headers,
    pathParams: req.pathParams,
    body: req.body,
  };
  res.write(JSON.stringify(out));
  res.end('');
}
`

func TestEngineEval_SetsStatusHeaderBody(t *testing.T) {
	e := &Engine{Script: sampleScript}
	req := &Request{
		Method: "POST",
		Path:   "/api/test",
		URL:    "http://example.com/api/test?x=1",
		Query: map[string]string{
			"x": "1",
		},
		Headers: map[string]string{
			"X-Req":        "abc",
			"Content-Type": "text/plain",
		},
		PathParams: map[string]string{"id": "42"},
		Body:       []byte("hello"),
	}
	resp, err := e.Eval(req)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if resp.Status != 201 {
		t.Fatalf("unexpected status: %d", resp.Status)
	}
	if ct := resp.Headers["Content-Type"]; ct != "application/json" {
		t.Fatalf("unexpected content type: %q", ct)
	}
	if !resp.Ended {
		t.Fatalf("expected response to be ended")
	}
	body := resp.Body.String()
	if !strings.Contains(body, `"method":"POST"`) {
		t.Fatalf("body missing method: %s", body)
	}
	if !strings.Contains(body, `"path":"/api/test"`) {
		t.Fatalf("body missing path: %s", body)
	}
	if !strings.Contains(body, `"body":"hello"`) {
		t.Fatalf("body missing body content: %s", body)
	}
}

func TestEngineEval_NoScriptError(t *testing.T) {
	e := &Engine{}
	_, err := e.Eval(&Request{})
	if err == nil {
		t.Fatalf("expected error for empty script")
	}
}
