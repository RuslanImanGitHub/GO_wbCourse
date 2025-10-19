package calendar

import (
	"testing"
	"time"
)

func TestCalendarService(t *testing.T) {
	service := NewService()

	t.Run("CreateEvent", func(t *testing.T) {
		req := CreateEventRequest{
			UserID: 1,
			Date:   "2025-12-31",
			Title:  "New Year Party",
		}

		event, err := service.CreateEvent(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if event.UserID != req.UserID {
			t.Errorf("Expected UserID %d, got %d", req.UserID, event.UserID)
		}

		if event.Title != req.Title {
			t.Errorf("Expected Title %s, got %s", req.Title, event.Title)
		}
	})

	t.Run("CreateEvent_InvalidDate", func(t *testing.T) {
		req := CreateEventRequest{
			UserID: 1,
			Date:   "invalid-date",
			Title:  "Test",
		}

		_, err := service.CreateEvent(req)
		if err != ErrInvalidDate {
			t.Errorf("Expected ErrInvalidDate, got %v", err)
		}
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		// First create an event
		createReq := CreateEventRequest{
			UserID: 1,
			Date:   "2025-12-31",
			Title:  "Original Title",
		}
		event, _ := service.CreateEvent(createReq)

		// Then update it
		updateReq := UpdateEventRequest{
			ID:    event.ID,
			Title: "Updated Title",
			Date:  "2026-01-01",
		}

		updated, err := service.UpdateEvent(updateReq)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if updated.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
		}
	})

	t.Run("GetEventsForDay", func(t *testing.T) {
		date := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		events, err := service.GetEventsForDay(1, date)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(events) == 0 {
			t.Error("Expected at least one event")
		}
	})
}
