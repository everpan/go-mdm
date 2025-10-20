package server

import (
	"bytes"
	"context"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// NewHertz creates a Hertz server that forwards all HTTP requests to the provided
// JavaScript script through the goja Engine. The script must define:
//
//	function handle(req, res) { /* ... */ }
//
// The server registers a catch-all route and delegates processing to the engine.
func NewHertz(script string) *server.Hertz {
	hs := server.New()
	eng := &Engine{Script: script}

	// Catch-all route for all methods and paths
	hs.Any("/*path", func(c context.Context, ctx *app.RequestContext) {
		// Build Engine Request from Hertz context
		req := &Request{
			Method:     string(ctx.Method()),
			Path:       string(ctx.Path()),
			URL:        ctx.Request.URI().String(),
			Query:      map[string]string{},
			Headers:    map[string]string{},
			PathParams: map[string]string{},
		}

		// Query params (only first value for simplicity)
		ctx.QueryArgs().VisitAll(func(k, v []byte) {
			req.Query[string(k)] = string(v)
		})

		// Headers (only first value)
		ctx.Request.Header.VisitAll(func(k, v []byte) {
			req.Headers[string(k)] = string(v)
		})

		// Path params
		if p := ctx.Params; len(p) > 0 {
			for _, kv := range p {
				req.PathParams[kv.Key] = kv.Value
			}
		}

		// Body
		if b := ctx.Request.BodyBytes(); len(b) > 0 {
			// Make a copy since ctx buffer can be reused
			req.Body = append([]byte(nil), b...)
		}

		// Decode URL (best-effort) so JS scripts can use a decoded string if needed
		if u, err := url.QueryUnescape(req.URL); err == nil {
			req.URL = u
		}

		resp, err := eng.Eval(req)
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")
			ctx.Response.SetBodyString(err.Error())
			return
		}

		// Write response back to Hertz
		if resp.Status != 0 {
			ctx.SetStatusCode(resp.Status)
		}
		for k, v := range resp.Headers {
			ctx.Response.Header.Set(k, v)
		}
		if resp.Body.Len() > 0 {
			// Copy to avoid referencing internal buffer beyond handler lifetime
			ctx.Response.SetBody(bytes.Clone(resp.Body.Bytes()))
		}
	})

	return hs
}
