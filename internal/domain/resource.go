package domain

import "time"

// Resource is anything that can be booked: a room, a court, a stylist...
type Resource struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
