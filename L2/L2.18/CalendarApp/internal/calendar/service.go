package calendar

import (
	"sync"
	"time"
)

type service struct {
	mu     sync.RWMutex
	events map[int]Event
	nextID int
}

func NewService() Service {
	return &service{
		events: make(map[int]Event),
		nextID: 1,
	}
}

// CreateEvent implements Service.
func (s *service) CreateEvent(req CreateEventRequest) (*Event, error) {
	if req.UserID == 0 || req.Title == "" || req.Date == "" {
		return nil, ErrInvalidInput
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event := Event{
		ID:     s.nextID,
		UserID: req.UserID,
		Date:   date,
		Title:  req.Title,
	}

	s.events[event.ID] = event
	s.nextID++

	return &event, nil
}

// UpdateEvent implements Service.
func (s *service) UpdateEvent(req UpdateEventRequest) (*Event, error) {
	if req.ID == 0 {
		return nil, ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event, exist := s.events[req.ID]
	if !exist {
		return nil, ErrEventNotFound
	}

	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, ErrInvalidDate
		}
		event.Date = date
	}

	if req.Title != "" {
		event.Title = req.Title
	}

	s.events[event.ID] = event
	return &event, nil
}

// DeleteEvent implements Service.
func (s *service) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.events[id]; !exist {
		return ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

// GetEventsForDay implements Service.
func (s *service) GetEventsForDay(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var events []Event
	for _, event := range s.events {
		if event.UserID == userID && isSameDay(event.Date, date) {
			events = append(events, event)
		}
	}
	return events, nil
}

// GetEventsForWeek implements Service.
func (s *service) GetEventsForWeek(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	year, week := date.ISOWeek()
	var events []Event
	for _, event := range s.events {
		if event.UserID == userID {
			eventYear, eventWeek := event.Date.ISOWeek()
			if eventWeek == week && eventYear == year {
				events = append(events, event)
			}
		}
	}

	return events, nil
}

// GetEventsForMonth implements Service.
func (s *service) GetEventsForMonth(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	year, month := date.Year(), date.Month()
	var events []Event
	for _, event := range s.events {
		if event.UserID == userID {
			eventYear, eventMonth := event.Date.Year(), event.Date.Month()
			if eventMonth == month && eventYear == year {
				events = append(events, event)
			}
		}
	}

	return events, nil
}

func isSameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
