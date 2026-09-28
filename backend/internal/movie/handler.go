package movie

import (
	"fmt"
	"net/http"
	"strconv"

	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/httpx"
)

// Handler serves the recommendation and movie detail endpoints.
type Handler struct {
	svc      *Service
	isOption map[string]bool
}

// NewHandler returns a Handler. options is the allowlist for option_id.
func NewHandler(svc *Service, options []domain.Option) *Handler {
	isOption := make(map[string]bool, len(options))
	for _, o := range options {
		isOption[o.ID] = true
	}
	return &Handler{svc: svc, isOption: isOption}
}

// Register adds the movie routes to mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/options/{option_id}/recommendations", h.ListRecommendations)
	mux.HandleFunc("GET /api/v1/movies/{movie_id}", h.GetMovie)
}

// ListRecommendations handles GET /api/v1/options/{option_id}/recommendations.
func (h *Handler) ListRecommendations(w http.ResponseWriter, r *http.Request) {
	optionID := r.PathValue("option_id")
	details := map[string]any{"option_id": optionID}
	if !h.isOption[optionID] {
		msg := fmt.Sprintf("option %q does not exist", optionID)
		httpx.WriteError(w, r, httpx.WithDetails(domain.ErrOptionNotFound, msg, details))
		return
	}
	recs, err := h.svc.Recommendations(optionID)
	if err != nil {
		httpx.WriteError(w, r, httpx.WithDetails(err, "", details))
		return
	}
	httpx.WriteData(w, r, http.StatusOK, newMovieSummaries(recs.Movies), newRecommendationsMeta(recs))
}

// GetMovie handles GET /api/v1/movies/{movie_id}.
func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	raw := r.PathValue("movie_id")
	movieID, err := strconv.Atoi(raw)
	if err != nil || movieID <= 0 {
		httpx.WriteError(w, r, httpx.WithDetails(domain.ErrValidation,
			"movie_id must be a positive integer", map[string]any{"movie_id": raw}))
		return
	}
	movie, err := h.svc.Movie(movieID)
	if err != nil {
		httpx.WriteError(w, r, httpx.WithDetails(err, "", map[string]any{"movie_id": movieID}))
		return
	}
	httpx.WriteData(w, r, http.StatusOK, newMovieDetail(movie), detailMeta{FetchedAt: formatTime(movie.FetchedAt)})
}
