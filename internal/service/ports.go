package service

import (
	"context"
	"time"

	"booking-app/internal/domain"
)

// Interfaces live with the consumer (the service), so the service never
// imports the database package.

type ResourceRepository interface {
	Create(ctx context.Context, r *domain.Resource) error
	GetByID(ctx context.Context, id string) (*domain.Resource, error)
	List(ctx context.Context) ([]domain.Resource, error)
}

type BookingRepository interface {
	Create(ctx context.Context, b *domain.Booking) error
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
	List(ctx context.Context, resourceID string) ([]domain.Booking, error)
	HasOverlap(ctx context.Context, resourceID string, start, end time.Time) (bool, error)
	UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error
}
