package domain

import "time"

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID            string        `json:"id"`
	ResourceID    string        `json:"resource_id"`
	CustomerName  string        `json:"customer_name"`
	CustomerEmail string        `json:"customer_email"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	Status        BookingStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
}
