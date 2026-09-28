// Package movie serves ranked recommendations per option and movie details,
// both read from the in-memory cache snapshot.
package movie

import (
	"fmt"
	"time"

	"pasta_night/be/internal/domain"
)

// MovieReader is the data access the service needs; Repository implements it.
type MovieReader interface { //nolint:revive // PROJECT-RULES.md §1 names this interface MovieReader
	ListWithDetails(optionID string) (domain.MovieList, []domain.MovieDetail, error)
	Detail(movieID int) (domain.MovieDetail, error)
}

// Recommendations are an option's movies in rank order with freshness info.
type Recommendations struct {
	OptionID  string
	Movies    []domain.MovieDetail
	FetchedAt time.Time
	IsStale   bool
}

// Service builds recommendations and movie details from cached data.
type Service struct {
	repo       MovieReader
	staleAfter time.Duration
	now        func() time.Time
}

// NewService returns a Service. staleAfter is LIST_TTL + REFRESH_CHECK_INTERVAL:
// a list older than that missed a refresh that was due (DATABASE.md §5).
func NewService(repo MovieReader, staleAfter time.Duration) *Service {
	return &Service{repo: repo, staleAfter: staleAfter, now: time.Now}
}

// Recommendations returns the movies of a valid option in rank order.
// Stale data is still returned, flagged with IsStale.
func (s *Service) Recommendations(optionID string) (Recommendations, error) {
	list, details, err := s.repo.ListWithDetails(optionID)
	if err != nil {
		return Recommendations{}, fmt.Errorf("list recommendations for %s: %w", optionID, err)
	}
	return Recommendations{
		OptionID:  optionID,
		Movies:    details,
		FetchedAt: list.FetchedAt,
		IsStale:   s.now().Sub(list.FetchedAt) > s.staleAfter,
	}, nil
}

// Movie returns one movie from a current recommendation list.
func (s *Service) Movie(movieID int) (domain.MovieDetail, error) {
	d, err := s.repo.Detail(movieID)
	if err != nil {
		return domain.MovieDetail{}, fmt.Errorf("get movie %d: %w", movieID, err)
	}
	return d, nil
}
