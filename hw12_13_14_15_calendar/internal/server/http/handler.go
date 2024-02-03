package internalhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	logger Logger
	app    Application
}

type EventRequest struct {
	Title            string        `json:"title,omitempty"`
	TimeFrom         time.Time     `json:"timeFrom,omitempty"`
	TimeTo           time.Time     `json:"timeTo,omitempty"`
	Description      string        `json:"description,omitempty"`
	UserID           string        `json:"userId,omitempty"`
	NotificationTime time.Duration `json:"notificationTime,omitempty"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewHandler(logger Logger, app Application) http.Handler {
	h := &handler{
		logger: logger,
		app:    app,
	}

	r := mux.NewRouter()
	r.HandleFunc("/healthz", h.Healthz).Methods(http.MethodGet)
	r.HandleFunc("/ready", h.Ready).Methods(http.MethodGet)
	r.HandleFunc("/events/day", h.GetEventsByDay).Methods(http.MethodGet)
	r.HandleFunc("/events/week", h.GetEventsByWeek).Methods(http.MethodGet)
	r.HandleFunc("/events/month", h.GetEventsByMonth).Methods(http.MethodGet)
	r.HandleFunc("/events", h.CreateEvent).Methods(http.MethodPost)
	r.HandleFunc("/events/{id}", h.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/events/{id}", h.DeleteEvent).Methods(http.MethodDelete)

	return r
}

func (s *handler) Healthz(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "OK"); err != nil {
		return
	}
}

func (s *handler) Ready(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "OK"); err != nil {
		return
	}
}

func (s *handler) GetEventsByDay(w http.ResponseWriter, r *http.Request) {

	unixTime := r.URL.Query().Get("time")

	if unixTime == "" {
		writeResponseError(w, fmt.Errorf("parameter 'time' is missing from URL"), s.logger)
		return
	}

	unixInt, err := strconv.ParseInt(unixTime, 10, 64)
	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	timeStart := time.Unix(unixInt, 0)

	events := s.app.GetEventsByDay(r.Context(), timeStart)

	writeResponseSuccess(w, events, s.logger)
}

func (s *handler) GetEventsByWeek(w http.ResponseWriter, r *http.Request) {

	unixTime := r.URL.Query().Get("time")

	if unixTime == "" {
		writeResponseError(w, fmt.Errorf("parameter 'time' is missing from URL"), s.logger)
		return
	}

	unixInt, err := strconv.ParseInt(unixTime, 10, 64)
	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	timeStart := time.Unix(unixInt, 0)

	events := s.app.GetEventsByWeek(r.Context(), timeStart)

	writeResponseSuccess(w, events, s.logger)
}

func (s *handler) GetEventsByMonth(w http.ResponseWriter, r *http.Request) {

	unixTime := r.URL.Query().Get("time")

	if unixTime == "" {
		writeResponseError(w, fmt.Errorf("parameter 'time' is missing from URL"), s.logger)
		return
	}

	unixInt, err := strconv.ParseInt(unixTime, 10, 64)
	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	timeStart := time.Unix(unixInt, 0)

	events := s.app.GetEventsByMonth(r.Context(), timeStart)

	writeResponseSuccess(w, events, s.logger)
}

func (s *handler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	var eventData EventRequest
	err := json.NewDecoder(r.Body).Decode(&eventData)

	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	id, err := s.app.CreateEvent(
		r.Context(),
		eventData.Title,
		eventData.TimeFrom,
		eventData.TimeTo,
		&eventData.Description,
		eventData.UserID,
		&eventData.NotificationTime,
	)

	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	writeResponseSuccess(w, id, s.logger)
}

func (s *handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var eventData EventRequest
	err := json.NewDecoder(r.Body).Decode(&eventData)

	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	event := storage.Event{
		ID:               id,
		Title:            eventData.Title,
		TimeFrom:         eventData.TimeFrom,
		TimeTo:           eventData.TimeTo,
		Description:      &eventData.Description,
		UserID:           eventData.UserID,
		NotificationTime: &eventData.NotificationTime,
	}

	err = s.app.UpdateEvent(
		r.Context(),
		id,
		event,
	)

	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	writeResponseSuccess(w, id, s.logger)
}

func (s *handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		writeResponseError(w, fmt.Errorf("parameter 'id' is missing from URL"), s.logger)
		return
	}

	err := s.app.DeleteEvent(r.Context(), id)

	if err != nil {
		writeResponseError(w, err, s.logger)
		return
	}

	writeResponseSuccess(w, 1, s.logger)
}

func writeResponseSuccess(w http.ResponseWriter, data interface{}, logger Logger) {

	response := &SuccessResponse{}
	response.Data = data

	buf, err := json.Marshal(response)
	if err != nil {
		logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(buf)
	if err != nil {
		logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func writeResponseError(w http.ResponseWriter, err error, logger Logger) {
	response := &ErrorResponse{}
	response.Error.Message = err.Error()

	buf, err := json.Marshal(response)
	if err != nil {
		logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(buf)
	if err != nil {
		logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
