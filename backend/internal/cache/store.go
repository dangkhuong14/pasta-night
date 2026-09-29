// Package cache holds the in-memory snapshot the read path serves from and
// persists it as JSON files under CACHE_DIR (DATABASE.md).
package cache

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"pasta_night/be/internal/domain"
)

// SchemaVersion is written into every cache file. Bump it whenever a
// persisted struct changes (DATABASE.md §6).
//
// v2: `providers` changed meaning from "providers in WATCH_REGION" to
// "providers merged across every region, most widely available first".
const SchemaVersion = 2

// Snapshot is an immutable view of the cache. Readers must never modify it
// or its maps; the Store swaps in a new Snapshot on every change.
type Snapshot struct {
	// Lists holds one list per option that has been refreshed, by option ID.
	Lists map[string]domain.MovieList
	// Details holds exactly the movies referenced by Lists, by movie ID.
	Details map[int]domain.MovieDetail
}

// Store owns the current Snapshot and the cache files. Reads are lock-free;
// writes (PublishList, GC) are serialized.
type Store struct {
	dir       string
	optionIDs []string
	isKnown   map[string]bool
	log       *slog.Logger

	current atomic.Pointer[Snapshot]
	mu      sync.Mutex // serializes writers so GC never runs between a publish's file writes
}

// NewStore returns a Store for the configured option IDs with an empty
// snapshot. Call Load to read existing files.
func NewStore(dir string, optionIDs []string, log *slog.Logger) *Store {
	s := &Store{
		dir:       dir,
		optionIDs: optionIDs,
		isKnown:   make(map[string]bool, len(optionIDs)),
		log:       log,
	}
	for _, id := range optionIDs {
		s.isKnown[id] = true
	}
	s.current.Store(&Snapshot{Lists: map[string]domain.MovieList{}, Details: map[int]domain.MovieDetail{}})
	return s
}

// Snapshot returns the current snapshot. It is never nil.
func (s *Store) Snapshot() *Snapshot {
	return s.current.Load()
}

// IsReady reports whether every configured option has a list.
func (s *Store) IsReady() bool {
	snap := s.current.Load()
	for _, id := range s.optionIDs {
		if _, ok := snap.Lists[id]; !ok {
			return false
		}
	}
	return true
}

// PublishList makes list and its details visible to readers. details must
// contain every movie in list.MovieIDs. Detail files are written before the
// list file so a list on disk never references a missing detail, and the
// snapshot is swapped last so readers never see a partial update.
func (s *Store) PublishList(ctx context.Context, list domain.MovieList, details []domain.MovieDetail) error {
	if !s.isKnown[list.OptionID] {
		return fmt.Errorf("publish list: unknown option %q", list.OptionID)
	}
	byID := make(map[int]domain.MovieDetail, len(details))
	for _, d := range details {
		byID[d.ID] = d
	}
	for _, id := range list.MovieIDs {
		if _, ok := byID[id]; !ok {
			return fmt.Errorf("publish list %s: no detail for movie %d", list.OptionID, id)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	cur := s.current.Load()
	for _, id := range list.MovieIDs {
		d := byID[id]
		// A detail reused from the current snapshot is already on disk. One
		// that dropped out of it since (and may be GC'd) is written again.
		if old, ok := cur.Details[id]; ok && old.FetchedAt.Equal(d.FetchedAt) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := s.writeDetail(d); err != nil {
			return fmt.Errorf("publish list %s: %w", list.OptionID, err)
		}
	}
	if err := s.writeList(list); err != nil {
		return fmt.Errorf("publish list %s: %w", list.OptionID, err)
	}

	lists := maps.Clone(cur.Lists)
	lists[list.OptionID] = list
	s.current.Store(&Snapshot{Lists: lists, Details: referencedDetails(lists, byID, cur.Details)})
	return nil
}

// referencedDetails returns the details the lists reference, preferring
// fresh over current. Unreferenced details drop out of memory here; GC
// deletes their files later.
func referencedDetails(lists map[string]domain.MovieList, fresh, current map[int]domain.MovieDetail) map[int]domain.MovieDetail {
	out := make(map[int]domain.MovieDetail)
	for _, list := range lists {
		for _, id := range list.MovieIDs {
			if d, ok := fresh[id]; ok {
				out[id] = d
			} else if d, ok := current[id]; ok {
				out[id] = d
			}
		}
	}
	return out
}

// GCStats counts the files GC deleted.
type GCStats struct {
	Details int
	Lists   int
}

// GC deletes detail files no list references and list files of options no
// longer configured. Call it after a refresh run has published its lists.
func (s *Store) GC(ctx context.Context) (GCStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap := s.current.Load()
	var stats GCStats
	var errs []error

	detailNames, err := jsonFileNames(filepath.Join(s.dir, detailsDir))
	if err != nil {
		return stats, fmt.Errorf("gc details: %w", err)
	}
	for _, name := range detailNames {
		if err := ctx.Err(); err != nil {
			return stats, err
		}
		id, ok := detailID(name)
		if !ok {
			continue // not a file this package wrote
		}
		if _, isReferenced := snap.Details[id]; isReferenced {
			continue
		}
		if err := removeFile(filepath.Join(s.dir, detailsDir, name)); err != nil {
			errs = append(errs, err)
			continue
		}
		stats.Details++
	}

	listNames, err := jsonFileNames(filepath.Join(s.dir, listsDir))
	if err != nil {
		return stats, fmt.Errorf("gc lists: %w", err)
	}
	for _, name := range listNames {
		if s.isKnown[strings.TrimSuffix(name, jsonExt)] {
			continue
		}
		if err := removeFile(filepath.Join(s.dir, listsDir, name)); err != nil {
			errs = append(errs, err)
			continue
		}
		stats.Lists++
	}
	return stats, errors.Join(errs...)
}

func removeFile(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}
