package cache

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"pasta_night/be/internal/domain"
)

var fetchedAt = time.Date(2026, 9, 27, 10, 0, 0, 0, time.UTC)

func newTestStore(t *testing.T, dir string) *Store {
	t.Helper()
	return NewStore(dir, []string{"friends", "solo"}, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func detail(id int) domain.MovieDetail {
	poster := fmt.Sprintf("https://image.tmdb.org/t/p/w500/p%d.jpg", id)
	year := 1999
	return domain.MovieDetail{
		ID:            id,
		Title:         fmt.Sprintf("Phim %d", id),
		OriginalTitle: fmt.Sprintf("Movie %d", id),
		PosterURL:     &poster,
		ReleaseYear:   &year,
		Rating:        8.2,
		Genres:        []string{"Phim Hành Động"},
		Directors:     []string{},
		Cast:          []string{"Keanu Reeves"},
		Providers:     []domain.Provider{{ID: 8, Name: "Netflix", Type: domain.ProviderFlatrate}},
		FetchedAt:     fetchedAt,
	}
}

func details(ids ...int) []domain.MovieDetail {
	out := make([]domain.MovieDetail, len(ids))
	for i, id := range ids {
		out[i] = detail(id)
	}
	return out
}

func list(optionID string, ids ...int) domain.MovieList {
	return domain.MovieList{OptionID: optionID, FetchedAt: fetchedAt, MovieIDs: ids}
}

func mustPublish(t *testing.T, s *Store, l domain.MovieList, d []domain.MovieDetail) {
	t.Helper()
	if err := s.PublishList(context.Background(), l, d); err != nil {
		t.Fatalf("PublishList(%s): %v", l.OptionID, err)
	}
}

func detailIDs(snap *Snapshot) []int {
	var ids []int
	for id := range snap.Details {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestPublishThenLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := newTestStore(t, dir)
	if _, err := s.Load(context.Background()); err != nil {
		t.Fatalf("Load empty dir: %v", err)
	}
	mustPublish(t, s, list("friends", 603, 680), details(603, 680))

	reloaded := newTestStore(t, dir)
	stats, err := reloaded.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if stats != (LoadStats{Lists: 1, Details: 2}) {
		t.Errorf("stats = %+v", stats)
	}
	if !reflect.DeepEqual(reloaded.Snapshot(), s.Snapshot()) {
		t.Errorf("reloaded snapshot differs:\n%+v\nwant\n%+v", reloaded.Snapshot(), s.Snapshot())
	}
}

func TestPublishFileFormat(t *testing.T) {
	dir := t.TempDir()
	s := newTestStore(t, dir)
	mustPublish(t, s, list("friends", 603), details(603))

	raw, err := os.ReadFile(filepath.Join(dir, "lists", "friends.json"))
	if err != nil {
		t.Fatal(err)
	}
	schemaLine := fmt.Sprintf(`"schema_version": %d`, SchemaVersion)
	for _, want := range []string{schemaLine, `"option_id": "friends"`, `"fetched_at": "2026-09-27T10:00:00Z"`} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("list file misses %s:\n%s", want, raw)
		}
	}
	raw, err = os.ReadFile(filepath.Join(dir, "details", "603.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{schemaLine, `"id": 603`, `"backdrop_url": null`, `"runtime_minutes": null`} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("detail file misses %s:\n%s", want, raw)
		}
	}
}

func TestPublishRejectsBadInput(t *testing.T) {
	tests := []struct {
		name    string
		list    domain.MovieList
		details []domain.MovieDetail
		wantErr string
	}{
		{"missing detail", list("friends", 1, 2), details(1), "no detail for movie 2"},
		{"unknown option", list("family", 1), details(1), "unknown option"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t, t.TempDir())
			err := s.PublishList(context.Background(), tt.list, tt.details)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want %q", err, tt.wantErr)
			}
			if len(s.Snapshot().Lists) != 0 {
				t.Error("snapshot changed after a rejected publish")
			}
		})
	}
}

func TestPublishKeepsOnlyReferencedDetailsInMemory(t *testing.T) {
	s := newTestStore(t, t.TempDir())
	mustPublish(t, s, list("friends", 1, 2), details(1, 2))
	mustPublish(t, s, list("solo", 2, 3), details(2, 3))
	if got := detailIDs(s.Snapshot()); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("details = %v, want [1 2 3]", got)
	}

	// Movie 2 stays: solo still lists it. Movie 1 drops out of memory.
	mustPublish(t, s, list("friends", 4), details(4))
	if got := detailIDs(s.Snapshot()); !reflect.DeepEqual(got, []int{2, 3, 4}) {
		t.Errorf("details = %v, want [2 3 4]", got)
	}
}

func TestGC(t *testing.T) {
	dir := t.TempDir()
	s := newTestStore(t, dir)
	mustPublish(t, s, list("friends", 1, 2), details(1, 2))
	mustPublish(t, s, list("friends", 2, 3), details(2, 3))
	writeFile(t, filepath.Join(dir, "lists", "removed-option.json"), `{"schema_version":1}`)
	writeFile(t, filepath.Join(dir, "details", "notes.json"), `{}`)

	stats, err := s.GC(context.Background())
	if err != nil {
		t.Fatalf("GC: %v", err)
	}
	if stats != (GCStats{Details: 1, Lists: 1}) {
		t.Errorf("stats = %+v, want 1 detail and 1 list", stats)
	}
	for path, want := range map[string]bool{
		"details/1.json":            false, // unreferenced
		"details/2.json":            true,
		"details/3.json":            true,
		"details/notes.json":        true, // not a file the store wrote
		"lists/friends.json":        true,
		"lists/removed-option.json": false, // option no longer configured
	} {
		if got := fileExists(filepath.Join(dir, path)); got != want {
			t.Errorf("%s exists = %v, want %v", path, got, want)
		}
	}
}

func TestPublishRewritesDetailDeletedByGC(t *testing.T) {
	dir := t.TempDir()
	s := newTestStore(t, dir)
	old := detail(1)
	mustPublish(t, s, list("friends", 1), []domain.MovieDetail{old})
	mustPublish(t, s, list("friends", 2), details(2))
	if _, err := s.GC(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fileExists(filepath.Join(dir, "details", "1.json")) {
		t.Fatal("GC should have deleted details/1.json")
	}

	// A refresh that reused the old in-memory detail must write it again.
	mustPublish(t, s, list("solo", 1), []domain.MovieDetail{old})
	if !fileExists(filepath.Join(dir, "details", "1.json")) {
		t.Error("details/1.json was not rewritten")
	}
}

func TestLoadSkipsInvalidFiles(t *testing.T) {
	dir := t.TempDir()
	valid := newTestStore(t, dir)
	mustPublish(t, valid, list("friends", 1), details(1))

	// Only detail 2 is skipped for its version; the others must fail their own
	// rule, so they carry the current SchemaVersion.
	writeFile(t, filepath.Join(dir, "details", "2.json"), `{"schema_version":99,"id":2}`)
	writeFile(t, filepath.Join(dir, "details", "3.json"), `{not json`)
	writeFile(t, filepath.Join(dir, "details", "4.json"), fmt.Sprintf(`{"schema_version":%d,"id":5}`, SchemaVersion))
	writeFile(t, filepath.Join(dir, "details", ".1.json.123.tmp"), `half-written`)
	writeFile(t, filepath.Join(dir, "lists", "solo.json"), fmt.Sprintf(
		`{"schema_version":%d,"option_id":"solo","fetched_at":"2026-09-27T10:00:00Z","movie_ids":[1,2]}`, SchemaVersion))
	writeFile(t, filepath.Join(dir, "lists", "removed-option.json"), fmt.Sprintf(
		`{"schema_version":%d,"option_id":"removed-option","fetched_at":"2026-09-27T10:00:00Z","movie_ids":[1]}`, SchemaVersion))

	s := newTestStore(t, dir)
	stats, err := s.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Skipped: details 2, 3, 4; list solo (movie 2 is invalid); list removed-option.
	if stats != (LoadStats{Lists: 1, Details: 1, Skipped: 5}) {
		t.Errorf("stats = %+v", stats)
	}
	if _, ok := s.Snapshot().Lists["solo"]; ok {
		t.Error("a list with a missing detail must be treated as missing")
	}
	if fileExists(filepath.Join(dir, "details", ".1.json.123.tmp")) {
		t.Error("temp file was not removed")
	}
}

func TestListFileValidate(t *testing.T) {
	available := map[int]domain.MovieDetail{1: detail(1), 2: detail(2)}
	manyIDs := make([]int, domain.MaxListSize+1)
	for i := range manyIDs {
		manyIDs[i] = i + 1
	}
	tests := []struct {
		name    string
		file    listFile
		wantErr string // empty = valid
	}{
		{"valid", listFile{SchemaVersion, list("friends", 1, 2)}, ""},
		{"wrong schema", listFile{SchemaVersion + 1, list("friends", 1)}, "schema_version"},
		{"option does not match file name", listFile{SchemaVersion, list("solo", 1)}, "file name"},
		{"no fetched_at", listFile{SchemaVersion, domain.MovieList{OptionID: "friends", MovieIDs: []int{1}}}, "fetched_at"},
		{"empty list", listFile{SchemaVersion, list("friends")}, "want 1-40"},
		{"too many movies", listFile{SchemaVersion, list("friends", manyIDs...)}, "want 1-40"},
		{"duplicate movie", listFile{SchemaVersion, list("friends", 1, 1)}, "twice"},
		{"missing detail", listFile{SchemaVersion, list("friends", 1, 3)}, "no valid detail"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.file.validate("friends", available)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validate: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestIsReady(t *testing.T) {
	s := newTestStore(t, t.TempDir())
	if s.IsReady() {
		t.Error("empty store must not be ready")
	}
	mustPublish(t, s, list("friends", 1), details(1))
	if s.IsReady() {
		t.Error("store must not be ready while solo has no list")
	}
	mustPublish(t, s, list("solo", 1), details(1))
	if !s.IsReady() {
		t.Error("store must be ready once every option has a list")
	}
}

func TestDetailID(t *testing.T) {
	tests := []struct {
		name   string
		want   int
		wantOK bool
	}{
		{"603.json", 603, true},
		{"0.json", 0, false},
		{"-1.json", 0, false},
		{"+603.json", 0, false},
		{"0603.json", 0, false},
		{"notes.json", 0, false},
	}
	for _, tt := range tests {
		got, ok := detailID(tt.name)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("detailID(%q) = %d, %v; want %d, %v", tt.name, got, ok, tt.want, tt.wantOK)
		}
	}
}
