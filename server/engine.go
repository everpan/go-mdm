package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
)

// Engine processes HTTP-like requests via a goja JavaScript script.
// It binds request/response objects for the script to read/write, then
// returns the produced response.
type Engine struct {
	// Script is the JavaScript code that will be executed for each request.
	// It should define a global function named `handle(req, res)`.
	Script string
}

// Request is a minimal HTTP request abstraction exposed to JS.
// Body is the raw bytes of the request body.
// Headers are canonicalized as-is (case-sensitive as provided by the caller).
// Query contains the parsed query parameters; values are strings (first value).
// PathParams contain path parameters when available.
// URL is the full request URL when available.
// Note: We keep it simple for unit-testability and portability across frameworks.
// The Hertz handler will adapt its context into this struct.
//
// All fields are exported so they can be marshaled or inspected if needed.
type Request struct {
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	URL        string            `json:"url"`
	Query      map[string]string `json:"query"`
	Headers    map[string]string `json:"headers"`
	PathParams map[string]string `json:"pathParams"`
	Body       []byte            `json:"body"`
}

// Response represents the outcome produced by JS.
// Body is written via res.write()/res.end().
// Status defaults to 200 unless changed.
// Headers can be set via res.setHeader(name, value).
// Ended indicates whether res.end() was called.
type Response struct {
	Status  int
	Headers map[string]string
	Body    bytes.Buffer
	Ended   bool
}

// Eval executes the engine's script and invokes handle(req, res).
func (e *Engine) Eval(req *Request) (*Response, error) {
	if e.Script == "" {
		return nil, fmt.Errorf("no script provided: Engine.Script is empty")
	}

	rt := goja.New()

	// Bind console.log for debugging convenience in scripts
	console := make(map[string]func(goja.FunctionCall) goja.Value)
	console["log"] = func(call goja.FunctionCall) goja.Value {
		args := make([]any, 0, len(call.Arguments))
		for _, a := range call.Arguments {
			args = append(args, a.Export())
		}
		// Best-effort JSON format for complex objects
		for i, v := range args {
			b, err := json.Marshal(v)
			if err == nil {
				args[i] = string(b)
			}
		}
		fmt.Println(args...)
		return goja.Undefined()
	}
	err := rt.Set("console", console)
	if err != nil {
		return nil, err
	}

	// Prepare JS-visible request object
	jsReq := map[string]any{
		"method":     req.Method,
		"path":       req.Path,
		"url":        req.URL,
		"query":      req.Query,
		"headers":    req.Headers,
		"pathParams": req.PathParams,
		"body":       string(req.Body), // expose as string for convenience
	}
	if err := rt.Set("req", jsReq); err != nil {
		return nil, fmt.Errorf("bind req: %w", err)
	}

	// Prepare response and JS-facing helpers
	resp := &Response{Status: 200, Headers: map[string]string{}}

	resObj := map[string]any{}
	resObj["setHeader"] = func(name, value string) {
		resp.Headers[name] = value
	}
	resObj["setStatus"] = func(status int) {
		resp.Status = status
	}
	resObj["write"] = func(chunk string) {
		resp.Body.WriteString(chunk)
	}
	resObj["end"] = func(final string) {
		if final != "" {
			resp.Body.WriteString(final)
		}
		resp.Ended = true
	}
	if err := rt.Set("res", resObj); err != nil {
		return nil, fmt.Errorf("bind res: %w", err)
	}

	// Load script and ensure a handle exists
	if _, err := rt.RunString(e.Script); err != nil {
		return nil, fmt.Errorf("evaluate script: %w", err)
	}

	fn := rt.Get("handle")
	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("script must define function handle(req, res)")
	}

	if _, err := callable(goja.Undefined(), rt.Get("req"), rt.Get("res")); err != nil {
		return nil, fmt.Errorf("invoke handle: %w", err)
	}

	return resp, nil
}
