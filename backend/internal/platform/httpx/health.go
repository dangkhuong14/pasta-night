package httpx

import "net/http"

// HealthzPath is the liveness endpoint, outside the /api/v1 prefix.
const HealthzPath = "/healthz"

type healthBody struct {
	Status       string `json:"status"`
	IsCacheReady bool   `json:"cache_ready"`
}

// Healthz answers 200 while the process is alive (API_SPEC.md §5.5). It is
// not wrapped in the envelope and is never cached, so probes see live state.
func Healthz(isReady func() bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		WriteJSON(w, r, http.StatusOK, healthBody{Status: "ok", IsCacheReady: isReady()})
	})
}
