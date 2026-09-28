// Package httpx holds the HTTP plumbing shared by feature handlers: the
// response envelope, the error-code table, routing fallbacks, and middleware.
package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// envelope is the success body of every /api/v1 response (API_SPEC.md §4).
type envelope struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

// fallbackErrorBody is sent when a response cannot be encoded at all.
var fallbackErrorBody = []byte(`{"code":"INTERNAL_ERROR","message":"internal server error","details":null}` + "\n")

// WriteData writes a success envelope. A nil meta becomes {}. Successful GET
// responses are cacheable for five minutes (API_SPEC.md §3).
func WriteData(w http.ResponseWriter, r *http.Request, status int, data, meta any) {
	if meta == nil {
		meta = struct{}{}
	}
	isGet := r.Method == http.MethodGet || r.Method == http.MethodHead
	if isGet && status >= 200 && status < 300 {
		w.Header().Set("Cache-Control", "public, max-age=300")
	}
	WriteJSON(w, r, status, envelope{Data: data, Meta: meta})
}

// WriteJSON writes v as a JSON body without the envelope. Feature handlers
// use WriteData; this is for infrastructure endpoints such as /healthz.
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		recordError(r, fmt.Errorf("encode response: %w", err))
		writeBody(w, http.StatusInternalServerError, fallbackErrorBody)
		return
	}
	writeBody(w, status, buf.Bytes())
}

func writeBody(w http.ResponseWriter, status int, body []byte) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, _ = w.Write(body) // the client went away; there is nobody left to tell
}
