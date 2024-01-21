package scripts

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type EventCreateResponse struct {
	Data string `json:"data"`
}

type EventCreate struct {
	Title            string        `json:"title,omitempty"`
	TimeFrom         time.Time     `json:"timeFrom,omitempty"`
	TimeTo           time.Time     `json:"timeTo,omitempty"`
	Description      string        `json:"description,omitempty"`
	UserID           string        `json:"userId,omitempty"`
	NotificationTime time.Duration `json:"notificationTime,omitempty"`
}

type EventsResponse struct {
	Data []EventResponse `json:"data"`
}

type EventResponse struct {
	ID               string        `json:"ID,omitempty"`
	Title            string        `json:"Title,omitempty"`
	TimeFrom         time.Time     `json:"TimeFrom,omitempty"`
	TimeTo           time.Time     `json:"TimeTo,omitempty"`
	Description      string        `json:"Description,omitempty"`
	UserID           string        `json:"UserId,omitempty"`
	NotificationTime time.Duration `json:"NotificationTime,omitempty"`
}

type EventDeleteResponse struct {
	Data int `json:"data"`
}

func TestIntegration(t *testing.T) {

	timeFrom := time.Date(2024, 1, 20, 14, 4, 5, 0, time.UTC)
	unixTime := timeFrom.Unix()

	newEvent := EventCreate{
		Title:            "New Event",
		TimeFrom:         timeFrom,
		TimeTo:           timeFrom.Add(24 * time.Hour),
		Description:      "Event description",
		UserID:           "user1",
		NotificationTime: 60,
	}

	var eventId string

	t.Run("creating event", func(t *testing.T) {

		payload, err := json.Marshal(newEvent)
		require.NoError(t, err)

		statusCode, result, err := sendRequest(http.MethodPost, "events", payload)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)

		var response EventCreateResponse

		err = json.Unmarshal(result, &response)
		require.NoError(t, err)
		require.NotEmpty(t, response.Data)

		eventId = response.Data
	})

	t.Run("getting event by day", func(t *testing.T) {

		statusCode, result, err := sendRequest(http.MethodGet, "events/day?time="+strconv.FormatInt(unixTime, 10), nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)

		var response EventsResponse

		err = json.Unmarshal(result, &response)
		require.NoError(t, err)
		require.Equal(t, 1, len(response.Data))
		require.Equal(t, eventId, response.Data[0].ID)
	})

	t.Run("getting event by week", func(t *testing.T) {

		statusCode, result, err := sendRequest(http.MethodGet, "events/week?time="+strconv.FormatInt(unixTime, 10), nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)

		var response EventsResponse

		err = json.Unmarshal(result, &response)
		require.NoError(t, err)
		require.Equal(t, 1, len(response.Data))
		require.Equal(t, eventId, response.Data[0].ID)
	})

	t.Run("getting event by month", func(t *testing.T) {

		statusCode, result, err := sendRequest(http.MethodGet, "events/month?time="+strconv.FormatInt(unixTime, 10), nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)

		var response EventsResponse

		err = json.Unmarshal(result, &response)
		require.NoError(t, err)
		require.Equal(t, 1, len(response.Data))
		require.Equal(t, eventId, response.Data[0].ID)
	})

	t.Run("getting event by day (POST METHOD)", func(t *testing.T) {

		statusCode, _, err := sendRequest(http.MethodPost, "events/day?time="+strconv.FormatInt(unixTime, 10), nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	t.Run("deleting event", func(t *testing.T) {

		statusCode, result, err := sendRequest(http.MethodDelete, "events/"+eventId, nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)

		var response EventDeleteResponse

		err = json.Unmarshal(result, &response)
		require.NoError(t, err)
		require.Equal(t, 1, response.Data)

		// Checking not exists event
		statusCode, result, err = sendRequest(http.MethodGet, "events/day?time="+strconv.FormatInt(unixTime, 10), nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)

		var eventsResponse EventsResponse

		err = json.Unmarshal(result, &response)
		require.NoError(t, err)
		require.Equal(t, 0, len(eventsResponse.Data))
	})

}

func sendRequest(method string, endpoint string, payload []byte) (int, []byte, error) {

	host := "http://calendar:8080/"
	//host := "http://localhost:8888/"

	req, err := http.NewRequest(method, host+endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, result, nil
}
