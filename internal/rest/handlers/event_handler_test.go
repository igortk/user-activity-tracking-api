package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"user-activity-tracking-api/internal/models"
)

type mockEventsRepository struct {
	createEventFn                 func(ctx context.Context, event *models.Event) error
	getEventsByUserIdAndDateRange func(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error)
}

func (m *mockEventsRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	if m.createEventFn != nil {
		return m.createEventFn(ctx, event)
	}
	return nil
}

func (m *mockEventsRepository) GetEventsByUserIdAndDateRange(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error) {
	if m.getEventsByUserIdAndDateRange != nil {
		return m.getEventsByUserIdAndDateRange(ctx, userID, limit, offset, from, to)
	}
	return nil, nil
}

func TestCreateActivityEvent_InvalidJSON(t *testing.T) {
	repo := &mockEventsRepository{}
	router := gin.Default()
	router.POST("/events", CreateActivityEvent(repo))

	body := bytes.NewBufferString(`{invalid json}`)

	req, _ := http.NewRequest(http.MethodPost, "/events", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateActivityEvent_ValidationFailed(t *testing.T) {
	repo := &mockEventsRepository{}
	router := gin.Default()
	router.POST("/events", CreateActivityEvent(repo))

	event := models.Event{}
	data, _ := json.Marshal(event)

	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateActivityEvent_RepoError(t *testing.T) {
	repo := &mockEventsRepository{
		createEventFn: func(ctx context.Context, event *models.Event) error {
			return errors.New("db error")
		},
	}

	router := gin.Default()
	router.POST("/events", CreateActivityEvent(repo))

	event := models.Event{
		UserID:               1,
		EventActionTimestamp: time.Now().UTC(),
		Action:               "created",
		Metadata:             json.RawMessage(`{"ip":"127.0.0.1","browser":"Chrome"}`),
	}

	data, _ := json.Marshal(event)
	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateActivityEvent_Success(t *testing.T) {
	repo := &mockEventsRepository{
		createEventFn: func(ctx context.Context, event *models.Event) error {
			return nil
		},
	}

	router := gin.Default()
	router.POST("/events", CreateActivityEvent(repo))

	event := models.Event{
		UserID:               1,
		EventActionTimestamp: time.Now().UTC(),
		Action:               "created",
		Metadata:             json.RawMessage(`{"ip":"127.0.0.1","browser":"Chrome"}`),
	}

	data, _ := json.Marshal(event)
	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if _, ok := resp["createdAt"]; !ok {
		t.Fatal("expected createdAt field in response")
	}
}

func TestGetActivityEventByUserIdDateRange_BindError(t *testing.T) {
	repo := &mockEventsRepository{}
	router := gin.Default()
	router.GET("/events", GetActivityEventByUserIdDateRange(repo))

	req, _ := http.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestGetActivityEventByUserIdDateRange_ValidationFailed(t *testing.T) {
	repo := &mockEventsRepository{}
	router := gin.Default()
	router.GET("/events", GetActivityEventByUserIdDateRange(repo))

	params := url.Values{}
	params.Add("user_id", "0")
	params.Add("from", time.Now().Add(-time.Hour).Format(time.RFC3339))
	params.Add("to", time.Now().Format(time.RFC3339))
	params.Add("offset", "0")
	params.Add("limit", "10")

	req, _ := http.NewRequest(http.MethodGet, "/events?"+params.Encode(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestGetActivityEventByUserIdDateRange_RepoError(t *testing.T) {
	repo := &mockEventsRepository{
		getEventsByUserIdAndDateRange: func(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error) {
			return nil, errors.New("db failure")
		},
	}

	router := gin.Default()
	router.GET("/events", GetActivityEventByUserIdDateRange(repo))

	params := url.Values{}
	params.Add("user_id", "1")
	params.Add("from", time.Now().Add(-time.Hour).Format(time.RFC3339))
	params.Add("to", time.Now().Format(time.RFC3339))
	params.Add("offset", "0")
	params.Add("limit", "10")

	req, _ := http.NewRequest(http.MethodGet, "/events?"+params.Encode(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

func TestGetActivityEventByUserIdDateRange_Success(t *testing.T) {
	expectedEvents := []models.Event{
		{
			UserID:               1,
			EventActionTimestamp: time.Now().Add(-5 * time.Minute).UTC(),
			Action:               "viewed",
			Metadata:             json.RawMessage(`{"ip":"127.0.0.1"}`),
		},
		{
			UserID:               1,
			EventActionTimestamp: time.Now().Add(-1 * time.Minute).UTC(),
			Action:               "updated",
			Metadata:             json.RawMessage(`{"page":"home"}`),
		},
	}

	repo := &mockEventsRepository{
		getEventsByUserIdAndDateRange: func(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error) {
			return expectedEvents, nil
		},
	}

	router := gin.Default()
	router.GET("/events", GetActivityEventByUserIdDateRange(repo))

	params := url.Values{}
	params.Add("user_id", "1")
	params.Add("from", time.Now().Add(-time.Hour).Format(time.RFC3339))
	params.Add("to", time.Now().Format(time.RFC3339))
	params.Add("offset", "0")
	params.Add("limit", "10")

	req, _ := http.NewRequest(http.MethodGet, "/events?"+params.Encode(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var got []models.Event
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(got) != len(expectedEvents) {
		t.Fatalf("expected %d events, got %d", len(expectedEvents), len(got))
	}
}
