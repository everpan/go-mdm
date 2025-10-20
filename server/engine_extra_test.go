package server

import (
	"strings"
	"testing"
)

func TestEngineEval_ScriptSyntaxError(t *testing.T) {
	e := &Engine{Script: `function handle(req, res) { res.setStatus(200); res.end('ok') `} // missing closing brace
	if _, err := e.Eval(&Request{}); err == nil {
		t.Fatalf("expected syntax error from script")
	}
}

func TestEngineEval_ConsoleLog_DoesNotPanic(t *testing.T) {
	e := &Engine{Script: `
	  console.log({a:1}, [1,2,3], 'x');
	  function handle(req, res){ res.end('x'); }
	`}
	resp, err := e.Eval(&Request{})
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if resp.Body.String() != "x" {
		t.Fatalf("unexpected body: %s", resp.Body.String())
	}
}

func TestEngineEval_Defaults_WhenNotSet(t *testing.T) {
	e := &Engine{Script: `function handle(req, res){ res.write('hello'); }`}
	resp, err := e.Eval(&Request{})
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if resp.Status != 200 {
		t.Fatalf("default status should be 200, got %d", resp.Status)
	}
	if resp.Ended {
		t.Fatalf("should not be ended")
	}
	if got := resp.Body.String(); got != "hello" {
		t.Fatalf("body %q", got)
	}
	if len(resp.Headers) != 0 {
		t.Fatalf("expected no headers, got %v", resp.Headers)
	}
}

func TestEngineEval_RequestFieldMappings_Minimal(t *testing.T) {
	// body contains some non-UTF8 bytes; since we cast to string, invalid bytes may be replaced.
	b := []byte{0xff, 0xfe, 'A'}
	req := &Request{
		Method:     "GET",
		Path:       "/p",
		URL:        "http://h/x",
		Query:      map[string]string{"a": "1"},
		Headers:    map[string]string{"K": "V"},
		PathParams: map[string]string{"id": "7"},
		Body:       b,
	}
	e := &Engine{Script: `
	  function handle(req, res){
	    if (req.method !== 'GET') throw new Error('m');
	    if (req.path !== '/p') throw new Error('p');
	    if (req.url.indexOf('http://h/x') !== 0) throw new Error('u');
	    if (req.query.a !== '1') throw new Error('q');
	    if (req.headers.K !== 'V') throw new Error('h');
	    if (req.pathParams.id !== '7') throw new Error('pp');
	    if (typeof req.body !== 'string') throw new Error('b');
	    res.end(req.body);
	  }
	`}
	resp, err := e.Eval(req)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if !strings.HasSuffix(resp.Body.String(), "A") { // last byte survives
		t.Fatalf("unexpected body mapping: %q", resp.Body.String())
	}
}
