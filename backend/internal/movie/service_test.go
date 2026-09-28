package movie

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"pasta_night/be/internal/domain"
)

// fakeReader serves one list and its details.
type fakeReader struct {
	list    domain.MovieList
	details []domain.MovieDetail
	err     error
}

func (f fakeReader) ListWithDetails(string) (domain.MovieList, []domain.MovieDetail, error) {
	return f.list, f.details, f.err
}

func (f fakeReader) Detail(movieID int) (domain.MovieDetail, error) {
	for _, d := range f.details {
		if d.ID == movieID {
			return d, nil
		}
	}
	return domain.MovieDetail{}, fmt.Errorf("movie %d: %w", movieID, domain.ErrMovieNotFound)
}

func TestRecommendationsStale(t *testing.T) {
	now := time.Date(2026, 9, 27, 12, 0, 0, 0, time.UTC)
	staleAfter := 169 * time.Hour // LIST_TTL 168h + REFRESH_CHECK_INTERVAL 1h
	tests := []struct {
		name      string
		age       time.Duration
		wantStale bool
	}{
		{"fresh", time.Hour, false},
		{"past LIST_TTL but refresh not yet overdue", 168*time.Hour + 30*time.Minute, false},
		{"exactly at the limit", staleAfter, false},
		{"refresh overdue", staleAfter + time.Second, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := fakeReader{
				list:    domain.MovieList{OptionID: "friends", FetchedAt: now.Add(-tt.age), MovieIDs: []int{1}},
				details: []domain.MovieDetail{{ID: 1, Title: "Phim 1"}},
			}
			svc := NewService(reader, staleAfter)
			svc.now = func() time.Time { return now }

			recs, err := svc.Recommendations("friends")
			if err != nil {
				t.Fatalf("Recommendations: %v", err)
			}
			if recs.IsStale != tt.wantStale {
				t.Errorf("IsStale = %v, want %v", recs.IsStale, tt.wantStale)
			}
			if recs.OptionID != "friends" || len(recs.Movies) != 1 || !recs.FetchedAt.Equal(reader.list.FetchedAt) {
				t.Errorf("recommendations = %+v", recs)
			}
		})
	}
}

func TestServiceErrors(t *testing.T) {
	svc := NewService(fakeReader{err: domain.ErrCacheNotReady}, time.Hour)
	if _, err := svc.Recommendations("friends"); !errors.Is(err, domain.ErrCacheNotReady) {
		t.Errorf("Recommendations err = %v, want ErrCacheNotReady", err)
	}
	if _, err := svc.Movie(42); !errors.Is(err, domain.ErrMovieNotFound) {
		t.Errorf("Movie err = %v, want ErrMovieNotFound", err)
	}
}
