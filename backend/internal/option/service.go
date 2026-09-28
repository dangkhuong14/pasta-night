// Package option serves the viewing options customers choose from.
package option

import "pasta_night/be/internal/domain"

// Service returns the configured viewing options.
type Service struct {
	options []domain.Option
}

// NewService returns a Service over options, in display order.
func NewService(options []domain.Option) *Service {
	return &Service{options: options}
}

// List returns every option in display order.
func (s *Service) List() []domain.Option {
	return s.options
}
