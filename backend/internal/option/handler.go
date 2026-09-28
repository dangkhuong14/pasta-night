package option

import (
	"net/http"

	"pasta_night/be/internal/platform/httpx"
)

// Handler serves the option endpoints.
type Handler struct {
	svc *Service
}

// NewHandler returns a Handler backed by svc.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register adds the option routes to mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/options", h.List)
}

// optionResponse is one element of GET /api/v1/options (API_SPEC.md §5.1).
// Discover params stay internal.
type optionResponse struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// List handles GET /api/v1/options.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	options := h.svc.List()
	data := make([]optionResponse, 0, len(options))
	for _, o := range options {
		data = append(data, optionResponse{ID: o.ID, Label: o.Label, Description: o.Description, Icon: o.Icon})
	}
	httpx.WriteData(w, r, http.StatusOK, data, nil)
}
