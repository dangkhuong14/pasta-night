package refresh

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
)

var testOptions = []domain.Option{
	{ID: "netflix-chill", Discover: friends.Discover},
	{ID: "solo", Discover: friends.Discover},
	{ID: "friends", Discover: friends.Discover},
}

func newTestJob(t *testing.T, client TMDB) (*Job, *cache.Store, string) {
	t.Helper()
	dir := t.TempDir()
	store := cache.NewStore(dir, []string{"netflix-chill", "solo", "friends"}, discardLogger())
	r := NewRefresher(client, store, Settings{WatchRegion: "VN", DiscoverPages: 1, Concurrency: 3, DetailTTL: 336 * time.Hour}, discardLogger())
	r.now = func() time.Time { return testNow }
	j := NewJob(r, store, testOptions, time.Hour, 168*time.Hour, discardLogger())
	j.now = func() time.Time { return testNow }
	return j, store, dir
}

// waitFor polls cond until it holds or the test times out.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestJobIsDue(t *testing.T) {
	tests := []struct {
		name    string
		age     time.Duration // 0 = no list
		wantDue bool
	}{
		{"no list", 0, true},
		{"fresh list", time.Hour, false},
		{"exactly LIST_TTL old", 168 * time.Hour, true},
		{"older than LIST_TTL", 200 * time.Hour, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, store, _ := newTestJob(t, newFakeTMDB())
			if tt.age > 0 {
				seedList(t, store, testNow.Add(-tt.age), 1)
			}
			if got := j.isDue("friends"); got != tt.wantDue {
				t.Errorf("isDue = %v, want %v", got, tt.wantDue)
			}
		})
	}
}

func TestJobTriggerQueuesInDisplayOrder(t *testing.T) {
	j, _, _ := newTestJob(t, newFakeTMDB())
	// Nothing runs the job: Trigger must still return at once.
	j.Trigger([]string{"friends"})
	j.Trigger([]string{"solo", "friends"})
	j.Trigger([]string{"solo"})

	if got := j.takePending(); !reflect.DeepEqual(got, []string{"solo", "friends"}) {
		t.Errorf("pending = %v, want [solo friends]", got)
	}
	if got := j.takePending(); len(got) != 0 {
		t.Errorf("pending after take = %v, want empty", got)
	}
}

func TestJobRunOnceSkipsFreshListsUnlessForced(t *testing.T) {
	fake := newFakeTMDB().withPage(1, 1, 1, 2).withMovies(1, 2)
	j, store, _ := newTestJob(t, fake)
	seedList(t, store, testNow.Add(-time.Hour), 1, 2)

	j.runOnce(context.Background(), []string{"friends"}, false)
	if calls, _ := fake.calls(); calls != 0 {
		t.Errorf("discover calls = %d for a fresh list, want 0", calls)
	}
	j.runOnce(context.Background(), []string{"friends"}, true)
	if calls, _ := fake.calls(); calls != 1 {
		t.Errorf("discover calls = %d for a forced refresh, want 1", calls)
	}
}

func TestJobRunOnceStopsOnUnauthorized(t *testing.T) {
	fake := newFakeTMDB()
	fake.discoverErr = errUnauthorized
	j, _, _ := newTestJob(t, fake)

	j.runOnce(context.Background(), []string{"netflix-chill", "solo", "friends"}, false)
	if calls, _ := fake.calls(); calls != 1 {
		t.Errorf("discover calls = %d, want 1: a rejected token stops the run", calls)
	}
}

func TestJobRun(t *testing.T) {
	fake := newFakeTMDB().withPage(1, 1, 1, 2, 3).withMovies(1, 2, 3)
	j, store, dir := newTestJob(t, fake)
	orphan := filepath.Join(dir, "details", "999.json")
	if err := os.MkdirAll(filepath.Dir(orphan), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(orphan, []byte(`{"schema_version":1,"id":999}`), 0o600); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		j.Run(ctx)
	}()
	defer func() {
		cancel()
		<-done
	}()

	// The first check runs immediately and refreshes every option.
	waitFor(t, "all options refreshed", store.IsReady)
	waitFor(t, "orphan detail deleted by GC", func() bool {
		_, err := os.Stat(orphan)
		return os.IsNotExist(err)
	})
	calls, _ := fake.calls()
	if calls != len(testOptions) {
		t.Errorf("discover calls = %d, want %d", calls, len(testOptions))
	}

	// A trigger forces a refresh even though every list is fresh.
	j.Trigger([]string{"friends"})
	waitFor(t, "forced refresh", func() bool {
		n, _ := fake.calls()
		return n == calls+1
	})

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after cancel")
	}
}
