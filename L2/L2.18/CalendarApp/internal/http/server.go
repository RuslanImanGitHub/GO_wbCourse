package http

import (
	"calendar/internal/calendar"
	"context"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port int, calendarService calendar.Service) *Server {
	handlers := NewHandler(calendarService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /create_event", handlers.CreateEvent)
	mux.HandleFunc("POST /update_event", handlers.UpdateEvent)
	mux.HandleFunc("POST /delete_event", handlers.DeleteEvent)
	mux.HandleFunc("GET /events_for_day", handlers.EventsForDay)
	mux.HandleFunc("GET /events_for_weekt", handlers.EventsForWeek)
	mux.HandleFunc("GET /events_for_month", handlers.EventsForMonth)

	// Логировщик через middleware
	handler := LoggingMiddleware(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
