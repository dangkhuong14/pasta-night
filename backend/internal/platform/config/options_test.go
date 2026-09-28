package config

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"pasta_night/be/internal/domain"
)

const soloYAML = `options:
  - id: solo
    label: "Một mình"
    description: "Thời gian cho riêng bạn"
    icon: user
    discover:
      genre_ids: [18, 53]
`

func TestParseOptionsAppliesDefaults(t *testing.T) {
	got, err := ParseOptions([]byte(soloYAML))
	if err != nil {
		t.Fatalf("ParseOptions: %v", err)
	}
	want := []domain.Option{{
		ID:          "solo",
		Label:       "Một mình",
		Description: "Thời gian cho riêng bạn",
		Icon:        "user",
		Discover: domain.DiscoverParams{
			GenreIDs:       []int{18, 53},
			GenreMode:      domain.GenreModeOr,
			MinVoteAverage: 6.5,
			MinVoteCount:   200,
			SortBy:         "popularity.desc",
		},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseOptions =\n%+v\nwant\n%+v", got, want)
	}
}

func TestParseOptionsExplicitZeroBeatsDefault(t *testing.T) {
	yml := strings.Replace(soloYAML, "genre_ids: [18, 53]",
		"genre_ids: [18]\n      genre_mode: and\n      min_vote_average: 0\n      min_vote_count: 0\n      sort_by: vote_average.desc", 1)
	got, err := ParseOptions([]byte(yml))
	if err != nil {
		t.Fatalf("ParseOptions: %v", err)
	}
	d := got[0].Discover
	if d.MinVoteAverage != 0 || d.MinVoteCount != 0 || d.GenreMode != domain.GenreModeAnd || d.SortBy != "vote_average.desc" {
		t.Errorf("discover = %+v, want explicit values kept", d)
	}
}

func TestParseOptionsErrors(t *testing.T) {
	replace := func(old, repl string) string { return strings.Replace(soloYAML, old, repl, 1) }
	withDiscover := func(extra string) string {
		return replace("genre_ids: [18, 53]", "genre_ids: [18]\n      "+extra)
	}
	duplicate := soloYAML + `  - id: solo
    label: "Khác"
    description: "Khác"
    icon: star
    discover:
      genre_ids: [35]
`
	tests := []struct {
		name    string
		yaml    string
		wantErr string
	}{
		{"empty file", "", "empty"},
		{"no options", "options: []\n", "no options"},
		{"malformed yaml", "options: [", "parse options yaml"},
		{"unknown field", replace("genre_ids:", "genres:"), "field genres not found"},
		{"id not kebab-case", replace("id: solo", "id: Solo_Mode"), "kebab-case"},
		{"missing label", replace(`label: "Một mình"`, `label: ""`), "label is required"},
		{"no genres", replace("genre_ids: [18, 53]", "genre_ids: []"), "genre_ids"},
		{"negative genre", replace("genre_ids: [18, 53]", "genre_ids: [18, -1]"), "positive"},
		{"negative provider", withDiscover("watch_provider_ids: [0]"), "watch_provider_ids"},
		{"bad genre mode", withDiscover("genre_mode: xor"), "genre_mode"},
		{"vote average above 10", withDiscover("min_vote_average: 11"), "min_vote_average"},
		{"negative vote count", withDiscover("min_vote_count: -1"), "min_vote_count"},
		{"bad sort_by", withDiscover("sort_by: popularity"), "sort_by"},
		{"duplicate id", duplicate, "duplicate id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseOptions([]byte(tt.yaml))
			if err == nil {
				t.Fatalf("ParseOptions succeeded, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q does not contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestShippedOptionsFileIsValid(t *testing.T) {
	options, err := LoadOptions(filepath.Join("..", "..", "..", "configs", "options.yaml"))
	if err != nil {
		t.Fatalf("LoadOptions: %v", err)
	}
	var ids []string
	for _, o := range options {
		ids = append(ids, o.ID)
	}
	if want := []string{"netflix-chill", "solo", "friends"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("option IDs = %v, want %v (display order)", ids, want)
	}
	netflix := options[0].Discover
	if netflix.MinVoteAverage != 7.0 || !reflect.DeepEqual(netflix.WatchProviderIDs, []int{8}) {
		t.Errorf("netflix-chill discover = %+v", netflix)
	}
}

func TestLoadOptionsMissingFile(t *testing.T) {
	if _, err := LoadOptions(filepath.Join(t.TempDir(), "missing.yaml")); err == nil {
		t.Fatal("want error for a missing file")
	}
}
