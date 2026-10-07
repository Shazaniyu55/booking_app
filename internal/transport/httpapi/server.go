package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/resources", func(r chi.Router) {
			r.Post("/", h.CreateResource)
			r.Get("/", h.ListResources)
			r.Get("/{id}", h.GetResource)
		})
		r.Route("/bookings", func(r chi.Router) {
			r.Post("/", h.CreateBooking)
			r.Get("/", h.ListBookings)
			r.Get("/{id}", h.GetBooking)
			r.Post("/{id}/cancel", h.CancelBooking)
		})
	})

	return r
}
