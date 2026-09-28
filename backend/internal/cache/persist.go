package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/fsutil"
)

// File layout under CACHE_DIR (DATABASE.md §6).
const (
	listsDir   = "lists"
	detailsDir = "details"
	jsonExt    = ".json"
	filePerm   = 0o644
	dirPerm    = 0o750
)

// listFile is the on-disk format of lists/{option_id}.json.
type listFile struct {
	SchemaVersion int `json:"schema_version"`
	domain.MovieList
}

// detailFile is the on-disk format of details/{movie_id}.json.
type detailFile struct {
	SchemaVersion int `json:"schema_version"`
	domain.MovieDetail
}

// LoadStats summarizes what Load kept.
type LoadStats struct {
	Lists   int
	Details int
	Skipped int // invalid files; refresh replaces them or GC deletes them
}

// Load creates the cache directories, removes temp files left by an
// interrupted write, and replaces the snapshot with the valid files on disk.
// Invalid files are logged and treated as missing (DATABASE.md §6); a list
// that references a missing or invalid detail is treated as missing too.
func (s *Store) Load(ctx context.Context) (LoadStats, error) {
	for _, sub := range []string{listsDir, detailsDir} {
		dir := filepath.Join(s.dir, sub)
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return LoadStats{}, fmt.Errorf("create cache dir: %w", err)
		}
		if err := removeTempFiles(dir); err != nil {
			return LoadStats{}, err
		}
	}
	details, skippedDetails, err := s.loadDetails(ctx)
	if err != nil {
		return LoadStats{}, err
	}
	lists, skippedLists, err := s.loadLists(ctx, details)
	if err != nil {
		return LoadStats{}, err
	}

	snap := &Snapshot{Lists: lists, Details: referencedDetails(lists, details, nil)}
	s.mu.Lock()
	s.current.Store(snap)
	s.mu.Unlock()
	return LoadStats{Lists: len(snap.Lists), Details: len(snap.Details), Skipped: skippedDetails + skippedLists}, nil
}

func (s *Store) loadDetails(ctx context.Context) (map[int]domain.MovieDetail, int, error) {
	dir := filepath.Join(s.dir, detailsDir)
	names, err := jsonFileNames(dir)
	if err != nil {
		return nil, 0, fmt.Errorf("list detail files: %w", err)
	}
	details := make(map[int]domain.MovieDetail, len(names))
	skipped := 0
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}
		var f detailFile
		if err := readJSON(filepath.Join(dir, name), &f); err != nil {
			s.skip(detailsDir, name, err)
			skipped++
			continue
		}
		if err := f.validate(name); err != nil {
			s.skip(detailsDir, name, err)
			skipped++
			continue
		}
		details[f.ID] = f.MovieDetail
	}
	return details, skipped, nil
}

func (f detailFile) validate(name string) error {
	if f.SchemaVersion != SchemaVersion {
		return fmt.Errorf("schema_version %d, want %d", f.SchemaVersion, SchemaVersion)
	}
	if id, ok := detailID(name); !ok || id != f.ID {
		return fmt.Errorf("id %d does not match the file name", f.ID)
	}
	return nil
}

func (s *Store) loadLists(ctx context.Context, details map[int]domain.MovieDetail) (map[string]domain.MovieList, int, error) {
	dir := filepath.Join(s.dir, listsDir)
	names, err := jsonFileNames(dir)
	if err != nil {
		return nil, 0, fmt.Errorf("list list files: %w", err)
	}
	lists := make(map[string]domain.MovieList, len(names))
	skipped := 0
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}
		optionID := strings.TrimSuffix(name, jsonExt)
		if !s.isKnown[optionID] {
			s.skip(listsDir, name, errors.New("option is not configured; GC deletes the file"))
			skipped++
			continue
		}
		var f listFile
		if err := readJSON(filepath.Join(dir, name), &f); err != nil {
			s.skip(listsDir, name, err)
			skipped++
			continue
		}
		if err := f.validate(optionID, details); err != nil {
			s.skip(listsDir, name, err)
			skipped++
			continue
		}
		lists[optionID] = f.MovieList
	}
	return lists, skipped, nil
}

func (f listFile) validate(optionID string, details map[int]domain.MovieDetail) error {
	switch {
	case f.SchemaVersion != SchemaVersion:
		return fmt.Errorf("schema_version %d, want %d", f.SchemaVersion, SchemaVersion)
	case f.OptionID != optionID:
		return fmt.Errorf("option_id %q does not match the file name", f.OptionID)
	case f.FetchedAt.IsZero():
		return errors.New("fetched_at is missing")
	case len(f.MovieIDs) == 0 || len(f.MovieIDs) > domain.MaxListSize:
		return fmt.Errorf("has %d movies, want 1-%d", len(f.MovieIDs), domain.MaxListSize)
	}
	seen := make(map[int]bool, len(f.MovieIDs))
	for _, id := range f.MovieIDs {
		if seen[id] {
			return fmt.Errorf("movie %d is listed twice", id)
		}
		seen[id] = true
		if _, ok := details[id]; !ok {
			return fmt.Errorf("movie %d has no valid detail file", id)
		}
	}
	return nil
}

func (s *Store) skip(sub, name string, err error) {
	s.log.Warn("ignoring cache file", "file", filepath.Join(sub, name), "err", err)
}

func (s *Store) writeDetail(d domain.MovieDetail) error {
	path := filepath.Join(s.dir, detailsDir, strconv.Itoa(d.ID)+jsonExt)
	return writeJSON(path, detailFile{SchemaVersion: SchemaVersion, MovieDetail: d})
}

func (s *Store) writeList(list domain.MovieList) error {
	path := filepath.Join(s.dir, listsDir, list.OptionID+jsonExt)
	return writeJSON(path, listFile{SchemaVersion: SchemaVersion, MovieList: list})
}

// writeJSON writes v atomically. It recreates a missing directory, so the
// cache heals if CACHE_DIR is wiped while the server runs.
func writeJSON(path string, v any) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}
	return fsutil.WriteFileAtomic(path, buf.Bytes(), filePerm)
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

// jsonFileNames returns the names of the regular *.json files in dir.
func jsonFileNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}
	var names []string
	for _, e := range entries {
		if e.Type().IsRegular() && strings.HasSuffix(e.Name(), jsonExt) {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// removeTempFiles deletes the temp files of writes interrupted by a crash.
func removeTempFiles(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", dir, err)
	}
	for _, e := range entries {
		if e.Type().IsRegular() && strings.HasSuffix(e.Name(), fsutil.TempSuffix) {
			if err := removeFile(filepath.Join(dir, e.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// detailID parses a detail file name such as "603.json".
func detailID(name string) (int, bool) {
	id, err := strconv.Atoi(strings.TrimSuffix(name, jsonExt))
	if err != nil || id <= 0 || strconv.Itoa(id)+jsonExt != name {
		return 0, false
	}
	return id, true
}
