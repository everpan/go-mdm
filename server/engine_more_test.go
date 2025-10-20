package server

import (
	"testing"
)

func TestEngineEval_ScriptWithoutHandle(t *testing.T) {
	e := &Engine{Script: `function notHandle() {}`}
	_, err := e.Eval(&Request{})
	if err == nil {
		t.Fatalf("expected error when handle is missing")
	}
}

func TestEngineEval_ScriptThrows(t *testing.T) {
	e := &Engine{Script: `function handle(req, res){ throw new Error('boom'); }`}
	_, err := e.Eval(&Request{})
	if err == nil {
		t.Fatalf("expected error when JS throws")
	}
}

func TestEngineEval_WriteWithoutEnd_AndMultipleWrites(t *testing.T) {
	e := &Engine{Script: `
	function handle(req, res) {
	  res.setStatus(202);
	  res.setHeader('X-A', '1');
	  res.write('hello');
	  res.write(' ');
	  res.write('world');
	  // don't call end here
	}
	`}
	resp, err := e.Eval(&Request{})
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if resp.Status != 202 {
		t.Fatalf("unexpected status: %d", resp.Status)
	}
	if resp.Headers["X-A"] != "1" {
		t.Fatalf("missing header")
	}
	if got := resp.Body.String(); got != "hello world" {
		t.Fatalf("unexpected body: %q", got)
	}
	if resp.Ended {
		t.Fatalf("should not be ended when end() wasn't called")
	}
}

func TestEngineEval_EndAppendsFinal(t *testing.T) {
	e := &Engine{Script: `
	function handle(req, res) {
	  res.write('A');
	  res.end('B');
	}
	`}
	resp, err := e.Eval(&Request{})
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if got := resp.Body.String(); got != "AB" {
		t.Fatalf("unexpected final body: %q", got)
	}
	if !resp.Ended {
		t.Fatalf("should be ended")
	}
}
