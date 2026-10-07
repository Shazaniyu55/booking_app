package httpapi

import (
	"net/http"
	"time"

	"booking-app/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	resources *service.ResourceService
	bookings  *service.BookingService
}

func NewHandler(r *service.ResourceService, b *service.BookingService) *Handler {
	return &Handler{resources: r, bookings: b}
}

// ---- resources ----

type createResourceReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req createResourceReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, err)
		return
	}
	res, err := h.resources.Create(r.Context(), req.Name, req.Description)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) ListResources(w http.ResponseWriter, r *http.Request) {
	out, err := h.resources.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
	res, err := h.resources.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// ---- bookings ----

type createBookingReq struct {
	ResourceID    string    `json:"resource_id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	StartTime     time.Time `json:"start_time"` // RFC3339, e.g. 2026-10-20T10:00:00Z
	EndTime       time.Time `json:"end_time"`
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req createBookingReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, err)
		return
	}
	b, err := h.bookings.Create(r.Context(), service.CreateBookingInput{
		ResourceID:    req.ResourceID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	out, err := h.bookings.List(r.Context(), r.URL.Query().Get("resource_id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	b, err := h.bookings.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	b, err := h.bookings.Cancel(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}
