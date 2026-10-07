package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"booking-app/internal/domain"
)

type ResourceService struct {
	repo ResourceRepository
}

func NewResourceService(repo ResourceRepository) *ResourceService {
	return &ResourceService{repo: repo}
}

func (s *ResourceService) Create(ctx context.Context, name, description string) (*domain.Resource, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	r := &domain.Resource{
		Name:        name,
		Description: strings.TrimSpace(description),
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *ResourceService) Get(ctx context.Context, id string) (*domain.Resource, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ResourceService) List(ctx context.Context) ([]domain.Resource, error) {
	return s.repo.List(ctx)
}
