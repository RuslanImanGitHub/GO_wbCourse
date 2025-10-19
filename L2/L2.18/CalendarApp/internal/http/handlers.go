package http

import (
	"calendar/internal/calendar"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type Handlers struct {
	calendarService calendar.Service
}

func NewHandler(calendarService calendar.Service) *Handlers {
	return &Handlers{
		calendarService: calendarService,
	}
}

func (h *Handlers) CreateEvent(respWriter http.ResponseWriter, r *http.Request) {
	var req calendar.CreateEventRequest
	if err := decodeRequest(r, &req); err != nil {
		writeError(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.calendarService.CreateEvent(req)
	if err != nil {
		writeCalendarError(respWriter, err)
		return
	}

	writeJSON(respWriter, Response{Result: event}, http.StatusOK)
}

func (h *Handlers) UpdateEvent(respWriter http.ResponseWriter, r *http.Request) {
	var req calendar.UpdateEventRequest
	if err := decodeRequest(r, &req); err != nil {
		writeError(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.calendarService.UpdateEvent(req)
	if err != nil {
		writeCalendarError(respWriter, err)
		return
	}

	writeJSON(respWriter, Response{Result: event}, http.StatusOK)
}

func (h *Handlers) DeleteEvent(respWriter http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(respWriter, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.calendarService.DeleteEvent(id); err != nil {
		writeCalendarError(respWriter, err)
		return
	}

	writeJSON(respWriter, Response{Result: "event deleted"}, http.StatusOK)
}

func (h *Handlers) EventsForDay(respWriter http.ResponseWriter, r *http.Request) {
	h.getEventsForPeriod(respWriter, r, h.calendarService.GetEventsForDay)
}

func (h *Handlers) EventsForWeek(respWriter http.ResponseWriter, r *http.Request) {
	h.getEventsForPeriod(respWriter, r, h.calendarService.GetEventsForWeek)
}

func (h *Handlers) EventsForMonth(respWriter http.ResponseWriter, r *http.Request) {
	h.getEventsForPeriod(respWriter, r, h.calendarService.GetEventsForMonth)
}

// #region EventGetter
type eventsGetter func(userID int, date time.Time) ([]calendar.Event, error)

func (h *Handlers) getEventsForPeriod(respWriter http.ResponseWriter, r *http.Request, getter eventsGetter) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		writeError(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := getter(userID, date)
	if err != nil {
		writeCalendarError(respWriter, err)
		return
	}

	writeJSON(respWriter, Response{Result: events}, http.StatusOK)
}

// #endregion

// #region HelperMethods
func decodeRequest(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" { //Для form data
		if err := r.ParseForm(); err != nil {
			return err
		}

		//Конвертим form в struct
		if req, ok := v.(*calendar.CreateEventRequest); ok {
			req.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
			req.Date = r.FormValue("date")
			req.Title = r.FormValue("title")
		} else if req, ok := v.(*calendar.UpdateEventRequest); ok {
			req.ID, _ = strconv.Atoi(r.FormValue("id"))
			req.Date = r.FormValue("date")
			req.Title = r.FormValue("title")
		}
		return nil
	}

	// По умолчанию пытаемся парсить как  JSON
	return json.NewDecoder(r.Body).Decode(v)
}

func writeError(respWriter http.ResponseWriter, errorMsg string, statusCode int) {
	writeJSON(respWriter, Response{Error: errorMsg}, statusCode)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeCalendarError(w http.ResponseWriter, err error) {
	switch err {
	case calendar.ErrEventNotFound:
		writeError(w, err.Error(), http.StatusServiceUnavailable)
	case calendar.ErrInvalidDate, calendar.ErrInvalidInput:
		writeError(w, err.Error(), http.StatusBadRequest)
	default:
		writeError(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseQueryParams(r *http.Request) (int, time.Time, error) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return 0, time.Time{}, calendar.ErrInvalidInput
	}

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, time.Time{}, calendar.ErrInvalidInput
	}

	return userID, date, nil
}

// #endregion
