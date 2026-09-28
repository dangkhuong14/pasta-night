package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

// assertJSON fails unless body is JSON equal to want.
func assertJSON(t *testing.T, body []byte, want string) {
	t.Helper()
	var got, exp any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("body is not JSON: %v\n%s", err, body)
	}
	if err := json.Unmarshal([]byte(want), &exp); err != nil {
		t.Fatalf("bad expected JSON: %v", err)
	}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("body =\n%s\nwant\n%s", body, want)
	}
}

// decodeError decodes an error response body.
func decodeError(t *testing.T, body []byte) errorBody {
	t.Helper()
	var e errorBody
	if err := json.Unmarshal(body, &e); err != nil {
		t.Fatalf("error body is not JSON: %v\n%s", err, body)
	}
	return e
}

// bufferLogger logs JSON lines, Debug and up, into the returned buffer.
func bufferLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})), &buf
}

func okHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	})
}
