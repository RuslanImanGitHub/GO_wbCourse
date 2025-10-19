package calendar

import (
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date")
	ErrInvalidInput  = errors.New("invalid input")
)

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
}

type CreateEventRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"` // YYYY-MM-DD
	Title  string `json:"title"`
}

type UpdateEventRequest struct {
	ID    int    `json:"user_id"`
	Date  string `json:"date,omitempty"` // YYYY-MM-DD
	Title string `json:"title,omitempty"`
}

type Service interface {
	CreateEvent(req CreateEventRequest) (*Event, error)
	UpdateEvent(req UpdateEventRequest) (*Event, error)
	DeleteEvent(id int) error
	GetEventsForDay(userID int, date time.Time) ([]Event, error)
	GetEventsForWeek(userID int, date time.Time) ([]Event, error)
	GetEventsForMonth(userID int, date time.Time) ([]Event, error)
}
