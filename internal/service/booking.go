package service

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"booking-app/internal/domain"
)

type BookingService struct {
	bookings  BookingRepository
	resources ResourceRepository
}

func NewBookingService(b BookingRepository, r ResourceRepository) *BookingService {
	return &BookingService{bookings: b, resources: r}
}

type CreateBookingInput struct {
	ResourceID    string
	CustomerName  string
	CustomerEmail string
	StartTime     time.Time
	EndTime       time.Time
}

func (s *BookingService) Create(ctx context.Context, in CreateBookingInput) (*domain.Booking, error) {
	in.CustomerName = strings.TrimSpace(in.CustomerName)
	if in.CustomerName == "" {
		return nil, fmt.Errorf("%w: customer_name is required", domain.ErrInvalid)
	}
	if _, err := mail.ParseAddress(in.CustomerEmail); err != nil {
		return nil, fmt.Errorf("%w: customer_email is not valid", domain.ErrInvalid)
	}
	if !in.EndTime.After(in.StartTime) {
		return nil, fmt.Errorf("%w: end_time must be after start_time", domain.ErrInvalid)
	}
	if in.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("%w: start_time must be in the future", domain.ErrInvalid)
	}

	// Resource must exist (returns ErrNotFound otherwise).
	if _, err := s.resources.GetByID(ctx, in.ResourceID); err != nil {
		return nil, err
	}

	taken, err := s.bookings.HasOverlap(ctx, in.ResourceID, in.StartTime, in.EndTime)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, fmt.Errorf("%w: time slot is already booked", domain.ErrConflict)
	}

	b := &domain.Booking{
		ResourceID:    in.ResourceID,
		CustomerName:  in.CustomerName,
		CustomerEmail: in.CustomerEmail,
		StartTime:     in.StartTime.UTC(),
		EndTime:       in.EndTime.UTC(),
		Status:        domain.StatusConfirmed,
		CreatedAt:     time.Now().UTC(),
	}
	if err := s.bookings.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BookingService) Get(ctx context.Context, id string) (*domain.Booking, error) {
	return s.bookings.GetByID(ctx, id)
}

func (s *BookingService) List(ctx context.Context, resourceID string) ([]domain.Booking, error) {
	return s.bookings.List(ctx, resourceID)
}

func (s *BookingService) Cancel(ctx context.Context, id string) (*domain.Booking, error) {
	b, err := s.bookings.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b.Status == domain.StatusCancelled {
		return nil, fmt.Errorf("%w: booking is already cancelled", domain.ErrConflict)
	}
	if err := s.bookings.UpdateStatus(ctx, id, domain.StatusCancelled); err != nil {
		return nil, err
	}
	b.Status = domain.StatusCancelled
	return b, nil
}
