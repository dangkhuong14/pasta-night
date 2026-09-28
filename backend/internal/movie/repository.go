package movie

import (
	"fmt"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
)

// Repository reads movie data from the cache snapshot. Each method reads a
// single snapshot, so a list and its details are always consistent.
type Repository struct {
	store *cache.Store
}

// NewRepository returns a Repository over store.
func NewRepository(store *cache.Store) *Repository {
	return &Repository{store: store}
}

// ListWithDetails returns an option's list and its details in rank order.
// It fails with domain.ErrCacheNotReady while the option has no list yet.
func (r *Repository) ListWithDetails(optionID string) (domain.MovieList, []domain.MovieDetail, error) {
	snap := r.store.Snapshot()
	list, ok := snap.Lists[optionID]
	if !ok {
		return domain.MovieList{}, nil, fmt.Errorf("list %s: %w", optionID, domain.ErrCacheNotReady)
	}
	details := make([]domain.MovieDetail, 0, len(list.MovieIDs))
	for _, id := range list.MovieIDs {
		if d, ok := snap.Details[id]; ok { // the store publishes details with their list
			details = append(details, d)
		}
	}
	return list, details, nil
}

// Detail returns a movie that is in a current list. The snapshot holds
// exactly the listed movies, so any other ID is domain.ErrMovieNotFound.
func (r *Repository) Detail(movieID int) (domain.MovieDetail, error) {
	d, ok := r.store.Snapshot().Details[movieID]
	if !ok {
		return domain.MovieDetail{}, fmt.Errorf("movie %d: %w", movieID, domain.ErrMovieNotFound)
	}
	return d, nil
}
