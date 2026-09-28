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
				Cast:           []string{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss", "Hugo Weaving", "Gloria Foster"},
				TrailerURL:     ptr("https://www.youtube.com/watch?v=viTrailer01"),
				Providers: []domain.Provider{
					{ID: 2, Name: "Apple TV", LogoURL: ptr("https://image.tmdb.org/t/p/w92/9ghgSC0MA082EL6HLCW3GalykFD.jpg"), Type: domain.ProviderRent},
					{ID: 8, Name: "Netflix", LogoURL: ptr("https://image.tmdb.org/t/p/w92/pbpMk2JmcoNnQwx5JGpXngfoWtp.jpg"), Type: domain.ProviderFlatrate},
					{ID: 3, Name: "Google Play Movies", LogoURL: nil, Type: domain.ProviderBuy},
				},
				FetchedAt: mappedAt,
			},
		},
		{
			fixture: "movie_no_region.json",
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
				Cast:           []string{"Leonardo DiCaprio"},
				TrailerURL:     nil, // only a teaser
				Providers:      []domain.Provider{},
				FetchedAt:      mappedAt,
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
				Cast:          []string{},
				Providers:     []domain.Provider{},
				FetchedAt:     mappedAt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			got, ok := mapMovieDetails(loadMovieFixture(t, tt.fixture), "VN", mappedAt)
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
			if _, ok := mapMovieDetails(tt.m, "VN", mappedAt); ok {
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
			got, ok := mapMovieDetails(tt.m, "VN", mappedAt)
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
		{"region missing", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"US": {Flatrate: []tmdb.Provider{p(1899, "Max", 1)}},
		}}, []string{}},
		{"no results at all", tmdb.WatchProviders{}, []string{}},
		{"first type wins, then display_priority", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"VN": {
				Flatrate: []tmdb.Provider{p(8, "Netflix", 3)},
				Free:     []tmdb.Provider{p(100, "Free TV", 1)},
				Rent:     []tmdb.Provider{p(8, "Netflix", 0)},
			},
		}}, []string{"100:free", "8:flatrate"}},
		{"equal priorities keep walk order", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"VN": {Ads: []tmdb.Provider{p(5, "Ads TV", 1)}, Buy: []tmdb.Provider{p(6, "Store", 1)}, Flatrate: []tmdb.Provider{p(7, "Stream", 1)}},
		}}, []string{"7:flatrate", "5:ads", "6:buy"}},
		{"invalid entries are skipped", tmdb.WatchProviders{Results: map[string]tmdb.RegionProviders{
			"VN": {Flatrate: []tmdb.Provider{p(0, "No ID", 1), p(9, "", 1), p(8, "Netflix", 2)}},
		}}, []string{"8:flatrate"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := providers(tt.wp, "VN")
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
		{Name: "third", Order: 2}, {Name: "first", Order: 0}, {Name: "sixth", Order: 5},
		{Name: "second", Order: 1}, {Name: "fourth", Order: 3}, {Name: "fifth", Order: 4},
	}
	if got := topCast(cast); !reflect.DeepEqual(got, []string{"first", "second", "third", "fourth", "fifth"}) {
		t.Errorf("cast = %v, want the first 5 in billing order", got)
	}
}
