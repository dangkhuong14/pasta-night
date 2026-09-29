package refresh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/tmdb"
)

var mappedAt = time.Date(2026, 9, 27, 10, 0, 5, 0, time.UTC)

func ptr[T any](v T) *T { return &v }

func loadMovieFixture(t *testing.T, name string) tmdb.MovieDetails {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "platform", "tmdb", "testdata", name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var m tmdb.MovieDetails
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	return m
}

func TestMapperFixtures(t *testing.T) {
	tests := []struct {
		fixture string
		want    domain.MovieDetail
	}{
		{
			fixture: "movie_603_full.json",
			want: domain.MovieDetail{
				ID:             603,
				Title:          "Ma Trận",
				OriginalTitle:  "The Matrix",
				Overview:       "Một hacker phát hiện thế giới anh đang sống chỉ là một mô phỏng.",
				Tagline:        "",
				PosterURL:      ptr("https://image.tmdb.org/t/p/w500/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg"),
				BackdropURL:    ptr("https://image.tmdb.org/t/p/w1280/fNG7i7RqMErkcqhohV2a6cV1Ehy.jpg"),
				ReleaseYear:    ptr(1999),
				Rating:         8.2,
				VoteCount:      26000,
				RuntimeMinutes: ptr(136),
				Genres:         []string{"Phim Hành Động", "Phim Khoa Học Viễn Tưởng"},
				Directors:      []string{"Lana Wachowski", "Lilly Wachowski"},
				// Gloria Foster has no profile_path in the fixture: her photo is nil.
				Cast: []domain.CastMember{
					{Name: "Keanu Reeves", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/8RZLOyYGsoRe9p44q3xin9QkMHv.jpg")},
					{Name: "Laurence Fishburne", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/2GbXERENPpl5MmlqOLlPVaVtifD.jpg")},
					{Name: "Carrie-Anne Moss", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/9zya72vRZYBQILfetACsnmCBgdj.jpg")},
					{Name: "Hugo Weaving", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/lSC8Et0PYi5zeQb3IpPkFje7hgR.jpg")},
					{Name: "Gloria Foster", ProfileURL: nil},
				},
				TrailerURL: ptr("https://www.youtube.com/watch?v=viTrailer01"),
				// Merged over US + VN: every provider is in one region, so the
				// type order decides, then display_priority (Max 1 < Netflix 5).
				Providers: []domain.Provider{
					{ID: 1899, Name: "Max", LogoURL: ptr("https://image.tmdb.org/t/p/w92/fksCUZ9QDWZMUwL2LgMtLckROUN.jpg"), Type: domain.ProviderFlatrate},
					{ID: 8, Name: "Netflix", LogoURL: ptr("https://image.tmdb.org/t/p/w92/pbpMk2JmcoNnQwx5JGpXngfoWtp.jpg"), Type: domain.ProviderFlatrate},
					{ID: 2, Name: "Apple TV", LogoURL: ptr("https://image.tmdb.org/t/p/w92/9ghgSC0MA082EL6HLCW3GalykFD.jpg"), Type: domain.ProviderRent},
					{ID: 3, Name: "Google Play Movies", LogoURL: nil, Type: domain.ProviderBuy},
				},
				FetchedAt: mappedAt,
			},
		},
		{
			fixture: "movie_other_region.json",
			want: domain.MovieDetail{
				ID:             27205,
				Title:          "Kẻ Đánh Cắp Giấc Mơ",
				OriginalTitle:  "Inception",
				Tagline:        "Tâm trí bạn là hiện trường vụ án.",
				PosterURL:      ptr("https://image.tmdb.org/t/p/w500/oYuLEt3zVCKq57qu2F8dT7NIa6f.jpg"),
				BackdropURL:    ptr("https://image.tmdb.org/t/p/w1280/8ZTVqvKDQ8emSGUEMjsS4yHAwrp.jpg"),
				ReleaseYear:    ptr(2010),
				Rating:         8.4,
				VoteCount:      36000,
				RuntimeMinutes: ptr(148),
				Genres:         []string{"Phim Hành Động", "Phim Khoa Học Viễn Tưởng", "Phim Phiêu Lưu"},
				Directors:      []string{},
				Cast: []domain.CastMember{
					{Name: "Leonardo DiCaprio", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/z2K1ERKlfMNzgLLIsM3jJ7Q9tHZ.jpg")},
				},
				TrailerURL: nil, // only a teaser
				// Offered only outside VN, and still listed: region no longer filters.
				Providers: []domain.Provider{
					{ID: 1899, Name: "Max", LogoURL: ptr("https://image.tmdb.org/t/p/w92/fksCUZ9QDWZMUwL2LgMtLckROUN.jpg"), Type: domain.ProviderFlatrate},
				},
				FetchedAt: mappedAt,
			},
		},
		{
			fixture: "movie_missing_fields.json",
			want: domain.MovieDetail{
				ID:            680,
				Title:         "Pulp Fiction", // falls back to original_title
				OriginalTitle: "Pulp Fiction",
				Genres:        []string{},
				Directors:     []string{},
				Cast:          []domain.CastMember{},
				Providers:     []domain.Provider{},
				FetchedAt:     mappedAt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			got, ok := mapMovieDetails(loadMovieFixture(t, tt.fixture), mappedAt)
			if !ok {
				t.Fatal("movie was dropped")
			}
			if !reflect.DeepEqual(got, tt.want) {
				gotJSON, _ := json.MarshalIndent(got, "", "  ")
				wantJSON, _ := json.MarshalIndent(tt.want, "", "  ")
				t.Errorf("mapped =\n%s\nwant\n%s", gotJSON, wantJSON)
			}
		})
	}
}

func TestMapperDrops(t *testing.T) {
	tests := []struct {
		name string
		m    tmdb.MovieDetails
	}{
		{"id 0", tmdb.MovieDetails{ID: 0, Title: "Phim"}},
		{"adult", tmdb.MovieDetails{ID: 1, Title: "Phim", Adult: true}},
		{"no title and no original title", tmdb.MovieDetails{ID: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := mapMovieDetails(tt.m, mappedAt); ok {
				t.Error("movie was kept, want dropped")
			}
		})
	}
}

func TestMapperScalarRules(t *testing.T) {
	tests := []struct {
		name string
		m    tmdb.MovieDetails
		want func(domain.MovieDetail) bool
	}{
		{"rating rounds to 1 decimal", tmdb.MovieDetails{ID: 1, Title: "x", VoteAverage: 8.216},
			func(d domain.MovieDetail) bool { return d.Rating == 8.2 }},
		{"rating rounds half up", tmdb.MovieDetails{ID: 1, Title: "x", VoteAverage: 7.25},
			func(d domain.MovieDetail) bool { return d.Rating == 7.3 }},
		{"runtime 0 is null", tmdb.MovieDetails{ID: 1, Title: "x", Runtime: 0},
			func(d domain.MovieDetail) bool { return d.RuntimeMinutes == nil }},
		{"release date empty is null", tmdb.MovieDetails{ID: 1, Title: "x", ReleaseDate: ""},
			func(d domain.MovieDetail) bool { return d.ReleaseYear == nil }},
		{"release date invalid is null", tmdb.MovieDetails{ID: 1, Title: "x", ReleaseDate: "abcd-01-01"},
			func(d domain.MovieDetail) bool { return d.ReleaseYear == nil }},
		{"release date too short is null", tmdb.MovieDetails{ID: 1, Title: "x", ReleaseDate: "199"},
			func(d domain.MovieDetail) bool { return d.ReleaseYear == nil }},
		{"empty poster path is null", tmdb.MovieDetails{ID: 1, Title: "x", PosterPath: ""},
			func(d domain.MovieDetail) bool { return d.PosterURL == nil && d.BackdropURL == nil }},
		{"localized title wins", tmdb.MovieDetails{ID: 1, Title: "Ma Trận", OriginalTitle: "The Matrix"},
			func(d domain.MovieDetail) bool { return d.Title == "Ma Trận" && d.OriginalTitle == "The Matrix" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := mapMovieDetails(tt.m, mappedAt)
			if !ok || !tt.want(got) {
				t.Errorf("mapped = %+v (ok %v)", got, ok)
			}
		})
	}
}

func TestMapperTrailerURL(t *testing.T) {
	yt := func(key, lang string, isOfficial bool) tmdb.Video {
		return tmdb.Video{Key: key, Site: "YouTube", Type: "Trailer", Language: lang, Official: isOfficial}
	}
	tests := []struct {
		name   string
		videos []tmdb.Video
		want   string // key; empty = null
	}{
		{"vi beats official en", []tmdb.Video{yt("en1", "en", true), yt("vi1", "vi", false)}, "vi1"},
		{"official first within a language", []tmdb.Video{yt("fan", "en", false), yt("official", "en", true)}, "official"},
		{"en beats other languages", []tmdb.Video{yt("fr", "fr", true), yt("en", "en", false)}, "en"},
		{"other language as a last resort", []tmdb.Video{yt("fr", "fr", false)}, "fr"},
		{"TMDB order breaks ties", []tmdb.Video{yt("first", "vi", true), yt("second", "vi", true)}, "first"},
		{"teasers and other sites are ignored", []tmdb.Video{
			{Key: "teaser", Site: "YouTube", Type: "Teaser", Language: "vi", Official: true},
			{Key: "vimeo", Site: "Vimeo", Type: "Trailer", Language: "vi", Official: true},
		}, ""},
		{"no videos", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trailerURL(tt.videos)
			switch {
			case tt.want == "" && got != nil:
				t.Errorf("trailer = %q, want null", *got)
			case tt.want != "" && (got == nil || *got != youtubeURL+tt.want):
				t.Errorf("trailer = %v, want %s", got, youtubeURL+tt.want)
			}
		})
	}
}

func TestMapperProviders(t *testing.T) {
	p := func(id int, name string, priority int) tmdb.Provider {
		return tmdb.Provider{ProviderID: id, ProviderName: name, LogoPath: "/l.jpg", DisplayPriority: priority}
	}
	tests := []struct {
		name string
		wp   tmdb.WatchProviders
		want []string // "id:type" in order
	}{
		{"no results at all", tmdb.WatchProviders{}, []string{}},
		{"a region we have no interest in still counts", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"US": {Flatrate: []tmdb.Provider{p(1899, "Max", 1)}},
		}}, []string{"1899:flatrate"}},
		{"more regions wins over display_priority", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"US": {Flatrate: []tmdb.Provider{p(8, "Netflix", 9)}, Rent: []tmdb.Provider{p(2, "Apple TV", 1)}},
			"DE": {Flatrate: []tmdb.Provider{p(8, "Netflix", 9)}},
			"FR": {Flatrate: []tmdb.Provider{p(8, "Netflix", 9)}},
		}}, []string{"8:flatrate", "2:rent"}},
		{"best type across regions wins", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"US": {Rent: []tmdb.Provider{p(8, "Netflix", 3)}},
			"DE": {Flatrate: []tmdb.Provider{p(8, "Netflix", 3)}},
		}}, []string{"8:flatrate"}},
		{"type order breaks a tie on region count", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"US": {Ads: []tmdb.Provider{p(5, "Ads TV", 1)}, Buy: []tmdb.Provider{p(6, "Store", 1)}, Flatrate: []tmdb.Provider{p(7, "Stream", 1)}},
		}}, []string{"7:flatrate", "5:ads", "6:buy"}},
		{"a provider listed twice in one region counts once", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"US": {Flatrate: []tmdb.Provider{p(8, "Netflix", 3)}, Buy: []tmdb.Provider{p(8, "Netflix", 3)}},
			"DE": {Flatrate: []tmdb.Provider{p(100, "Free TV", 1)}},
			"FR": {Flatrate: []tmdb.Provider{p(100, "Free TV", 1)}},
		}}, []string{"100:flatrate", "8:flatrate"}},
		{"invalid entries are skipped", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"VN": {Flatrate: []tmdb.Provider{p(0, "No ID", 1), p(9, "", 1), p(8, "Netflix", 2)}},
		}}, []string{"8:flatrate"}},
		{
			// TMDB gives one service several IDs: 9 and 119 are both
			// "Amazon Prime Video". Keep the better-ranked one only.
			name: "same name under different ids appears once",
			wp: tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
				"US": {Flatrate: []tmdb.Provider{p(9, "Amazon Prime Video", 3), p(8, "Netflix", 1)}},
				"DE": {Flatrate: []tmdb.Provider{p(119, "amazon prime video ", 2), p(8, "Netflix", 1)}},
				"FR": {Flatrate: []tmdb.Provider{p(119, "Amazon Prime Video", 2)}},
			},
			},
			want: []string{"8:flatrate", "119:flatrate"},
		},
		{"capped at maxProviders, most available first", manyRegionProviders(), []string{
			"1:flatrate", "2:flatrate", "3:flatrate", "4:flatrate",
			"5:flatrate", "6:flatrate", "7:flatrate", "8:flatrate",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := providers(tt.wp)
			if got == nil {
				t.Fatal("providers must be [] not nil")
			}
			keys := []string{}
			for _, pr := range got {
				keys = append(keys, fmtProvider(pr))
			}
			if !reflect.DeepEqual(keys, tt.want) {
				t.Errorf("providers = %v, want %v", keys, tt.want)
			}
		})
	}
}

func fmtProvider(p domain.Provider) string {
	return fmt.Sprintf("%d:%s", p.ID, p.Type)
}

// manyRegionProviders builds 12 providers offered by a decreasing number of
// regions, so provider 1 is the most widely available and 12 the least.
func manyRegionProviders() tmdb.WatchProviders {
	results := make(map[string]tmdb.RegionProviders)
	for id := 1; id <= 12; id++ {
		for region := 0; region <= 12-id; region++ {
			key := fmt.Sprintf("R%02d", region)
			rp := results[key]
			rp.Flatrate = append(rp.Flatrate, tmdb.Provider{
				ProviderID: id, ProviderName: fmt.Sprintf("P%d", id), LogoPath: "/l.jpg",
			})
			results[key] = rp
		}
	}
	return tmdb.WatchProviders{Results: results}
}

// TestMapperProvidersIsDeterministic guards against Go's random map iteration
// leaking into the cache files: the same input must always map the same way.
func TestMapperProvidersIsDeterministic(t *testing.T) {
	wp := manyRegionProviders()
	want := providers(wp)
	for i := 0; i < 50; i++ {
		if got := providers(wp); !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d = %v, want %v", i, got, want)
		}
	}
}

func TestMapperDirectorsAndCast(t *testing.T) {
	crew := []tmdb.CrewMember{
		{Name: "B", Job: "Director"},
		{Name: "A", Job: "Director"},
		{Name: "B", Job: "Director"},
		{Name: "C", Job: "Producer"},
		{Name: "", Job: "Director"},
	}
	if got := directors(crew); !reflect.DeepEqual(got, []string{"B", "A"}) {
		t.Errorf("directors = %v, want [B A] (TMDB order, de-duplicated)", got)
	}
	cast := []tmdb.CastMember{
		{Name: "third", Order: 2, ProfilePath: "/c.jpg"}, {Name: "first", Order: 0, ProfilePath: "/a.jpg"},
		{Name: "sixth", Order: 5, ProfilePath: "/f.jpg"}, {Name: "second", Order: 1, ProfilePath: ""},
		{Name: "fourth", Order: 3, ProfilePath: "/d.jpg"}, {Name: "fifth", Order: 4, ProfilePath: "/e.jpg"},
	}
	want := []domain.CastMember{
		{Name: "first", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/a.jpg")},
		{Name: "second", ProfileURL: nil}, // no profile_path → initials in the UI
		{Name: "third", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/c.jpg")},
		{Name: "fourth", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/d.jpg")},
		{Name: "fifth", ProfileURL: ptr("https://image.tmdb.org/t/p/w185/e.jpg")},
	}
	if got := topCast(cast); !reflect.DeepEqual(got, want) {
		t.Errorf("cast = %+v, want the first 5 in billing order with photos", got)
	}
}
