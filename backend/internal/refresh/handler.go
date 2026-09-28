package refresh

import (
	"fmt"
	"net/http"

	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/httpx"
)

// Handler serves the admin refresh endpoint.
type Handler struct {
	job       *Job
	optionIDs []string
	isOption  map[string]bool
}

// NewHandler returns a Handler that triggers job for the given options.
func NewHandler(job *Job, options []domain.Option) *Handler {
	h := &Handler{job: job, optionIDs: make([]string, 0, len(options)), isOption: make(map[string]bool, len(options))}
	for _, o := range options {
		h.optionIDs = append(h.optionIDs, o.ID)
		h.isOption[o.ID] = true
	}
	return h
}

// Register adds the admin routes to mux behind requireAdmin. Register them
// only when ADMIN_TOKEN is set: unregistered routes answer 404 NOT_FOUND,
// which is how disabled admin endpoints must look (API_SPEC.md §2).
func (h *Handler) Register(mux *http.ServeMux, requireAdmin func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/admin/refresh", requireAdmin(http.HandlerFunc(h.TriggerRefresh)))
}

type acceptedResponse struct {
	Accepted []string `json:"accepted"`
}

// TriggerRefresh handles POST /api/v1/admin/refresh. The refresh runs in the
// background; the response only lists the accepted options.
func (h *Handler) TriggerRefresh(w http.ResponseWriter, r *http.Request) {
	ids := h.optionIDs
	if optionID := r.URL.Query().Get("option_id"); optionID != "" {
		if !h.isOption[optionID] {
			msg := fmt.Sprintf("option %q does not exist", optionID)
			httpx.WriteError(w, r, httpx.WithDetails(domain.ErrOptionNotFound, msg, map[string]any{"option_id": optionID}))
			return
		}
		ids = []string{optionID}
	}
	h.job.Trigger(ids)
	httpx.WriteData(w, r, http.StatusAccepted, acceptedResponse{Accepted: ids}, nil)
}
