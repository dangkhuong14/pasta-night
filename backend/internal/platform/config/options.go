package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"pasta_night/be/internal/domain"
)

// Defaults for omitted discover fields (DATABASE.md §2 ViewingOption).
const (
	DefaultGenreMode      = domain.GenreModeOr
	DefaultMinVoteAverage = 6.5
	DefaultMinVoteCount   = 200
	DefaultSortBy         = "popularity.desc"
)

var (
	optionIDPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	sortByPattern   = regexp.MustCompile(`^[a-z_]+(\.[a-z_]+)*\.(asc|desc)$`)
)

// optionsFile mirrors configs/options.yaml. Pointer fields tell an omitted
// value from an explicit zero, so defaults apply only to omitted fields.
type optionsFile struct {
	Options []optionEntry `yaml:"options"`
}

type optionEntry struct {
	ID          string        `yaml:"id"`
	Label       string        `yaml:"label"`
	Description string        `yaml:"description"`
	Icon        string        `yaml:"icon"`
	Discover    discoverEntry `yaml:"discover"`
}

type discoverEntry struct {
	GenreIDs         []int    `yaml:"genre_ids"`
	GenreMode        string   `yaml:"genre_mode"`
	MinVoteAverage   *float64 `yaml:"min_vote_average"`
	MinVoteCount     *int     `yaml:"min_vote_count"`
	WatchProviderIDs []int    `yaml:"watch_provider_ids"`
	SortBy           string   `yaml:"sort_by"`
}

// LoadOptions reads and validates the viewing options file. The returned
// slice keeps file order, which is the display order.
func LoadOptions(path string) ([]domain.Option, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("read options file: %w", err)
	}
	return ParseOptions(data)
}

// ParseOptions parses and validates options YAML. Unknown fields are
// rejected so a typo cannot silently drop a filter.
func ParseOptions(data []byte) ([]domain.Option, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	var file optionsFile
	if err := dec.Decode(&file); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("options file is empty")
		}
		return nil, fmt.Errorf("parse options yaml: %w", err)
	}
	if len(file.Options) == 0 {
		return nil, errors.New("options file defines no options")
	}

	var errs []error
	seen := make(map[string]bool, len(file.Options))
	options := make([]domain.Option, 0, len(file.Options))
	for i, entry := range file.Options {
		opt, err := entry.toDomain()
		if err != nil {
			errs = append(errs, fmt.Errorf("options[%d] %q: %w", i, entry.ID, err))
			continue
		}
		if seen[opt.ID] {
			errs = append(errs, fmt.Errorf("options[%d]: duplicate id %q", i, opt.ID))
			continue
		}
		seen[opt.ID] = true
		options = append(options, opt)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return options, nil
}

func (e optionEntry) toDomain() (domain.Option, error) {
	var errs []error
	if !optionIDPattern.MatchString(e.ID) {
		errs = append(errs, fmt.Errorf("id %q must be kebab-case", e.ID))
	}
	requiredFields := []struct{ name, value string }{
		{"label", e.Label}, {"description", e.Description}, {"icon", e.Icon},
	}
	for _, f := range requiredFields {
		if strings.TrimSpace(f.value) == "" {
			errs = append(errs, fmt.Errorf("%s is required", f.name))
		}
	}
	discover, err := e.Discover.toDomain()
	if err != nil {
		errs = append(errs, err)
	}
	if err := errors.Join(errs...); err != nil {
		return domain.Option{}, err
	}
	return domain.Option{
		ID:          e.ID,
		Label:       e.Label,
		Description: e.Description,
		Icon:        e.Icon,
		Discover:    discover,
	}, nil
}

func (d discoverEntry) toDomain() (domain.DiscoverParams, error) {
	var errs []error
	if len(d.GenreIDs) == 0 {
		errs = append(errs, errors.New("discover.genre_ids needs at least one TMDB genre ID"))
	}
	errs = append(errs, positiveIDs("discover.genre_ids", d.GenreIDs), positiveIDs("discover.watch_provider_ids", d.WatchProviderIDs))

	mode := DefaultGenreMode
	if d.GenreMode != "" {
		mode = domain.GenreMode(d.GenreMode)
	}
	if mode != domain.GenreModeAnd && mode != domain.GenreModeOr {
		errs = append(errs, fmt.Errorf("discover.genre_mode must be %q or %q, got %q", domain.GenreModeAnd, domain.GenreModeOr, d.GenreMode))
	}

	minVoteAverage := DefaultMinVoteAverage
	if d.MinVoteAverage != nil {
		minVoteAverage = *d.MinVoteAverage
	}
	if minVoteAverage < 0 || minVoteAverage > 10 {
		errs = append(errs, fmt.Errorf("discover.min_vote_average must be between 0 and 10, got %v", minVoteAverage))
	}

	minVoteCount := DefaultMinVoteCount
	if d.MinVoteCount != nil {
		minVoteCount = *d.MinVoteCount
	}
	if minVoteCount < 0 {
		errs = append(errs, fmt.Errorf("discover.min_vote_count must be >= 0, got %d", minVoteCount))
	}

	sortBy := DefaultSortBy
	if d.SortBy != "" {
		sortBy = d.SortBy
	}
	if !sortByPattern.MatchString(sortBy) {
		errs = append(errs, fmt.Errorf("discover.sort_by must look like popularity.desc, got %q", sortBy))
	}

	if err := errors.Join(errs...); err != nil {
		return domain.DiscoverParams{}, err
	}
	return domain.DiscoverParams{
		GenreIDs:         d.GenreIDs,
		GenreMode:        mode,
		MinVoteAverage:   minVoteAverage,
		MinVoteCount:     minVoteCount,
		WatchProviderIDs: d.WatchProviderIDs,
		SortBy:           sortBy,
	}, nil
}

// positiveIDs returns an error if any TMDB ID is not positive, or nil.
func positiveIDs(field string, ids []int) error {
	for _, id := range ids {
		if id <= 0 {
			return fmt.Errorf("%s must contain positive IDs, got %d", field, id)
		}
	}
	return nil
}
