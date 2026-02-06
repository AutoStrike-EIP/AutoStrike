package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// withAuthAnalytics wraps a handler with authentication context for testing
func withAuthAnalytics(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		handler(c)
	}
}

func TestNewAnalyticsHandler(t *testing.T) {
	service := application.NewAnalyticsService(&mockResultRepoForHandler{})
	handler := NewAnalyticsHandler(service)

	if handler == nil {
		t.Fatal("NewAnalyticsHandler returned nil")
	}
}

func TestAnalyticsHandler_RegisterRoutes(t *testing.T) {
	service := application.NewAnalyticsService(&mockResultRepoForHandler{})
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	routes := router.Routes()
	expectedPaths := map[string]string{
		"/api/v1/analytics/compare": "GET",
		"/api/v1/analytics/trend":   "GET",
		"/api/v1/analytics/summary": "GET",
		"/api/v1/analytics/period":  "GET",
	}

	for path, method := range expectedPaths {
		found := false
		for _, route := range routes {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Route %s %s not found", method, path)
		}
	}
}

func TestAnalyticsHandler_CompareScores(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(10 * time.Minute)
	repo := &mockResultRepoForHandler{
		executions: []*entity.Execution{
			{
				ID:          "exec-1",
				ScenarioID:  "scenario-1",
				Status:      entity.ExecutionCompleted,
				StartedAt:   now.AddDate(0, 0, -2),
				CompletedAt: &completedAt,
				Score: &entity.SecurityScore{
					Overall:    80.0,
					Blocked:    4,
					Detected:   2,
					Successful: 4,
					Total:      10,
				},
			},
		},
	}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response application.ScoreComparison
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Current.ExecutionCount != 1 {
		t.Errorf("Current.ExecutionCount = %d, want 1", response.Current.ExecutionCount)
	}
}

func TestAnalyticsHandler_CompareScores_WithDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=30", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAnalyticsHandler_GetScoreTrend(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response application.ScoreTrend
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Period != "30d" {
		t.Errorf("Period = %q, want '30d'", response.Period)
	}
}

func TestAnalyticsHandler_GetScoreTrend_WithDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=7", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response application.ScoreTrend
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Period != "7d" {
		t.Errorf("Period = %q, want '7d'", response.Period)
	}
}

func TestAnalyticsHandler_GetExecutionSummary(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(10 * time.Minute)
	repo := &mockResultRepoForHandler{
		executions: []*entity.Execution{
			{
				ID:          "exec-1",
				ScenarioID:  "scenario-1",
				Status:      entity.ExecutionCompleted,
				StartedAt:   now.AddDate(0, 0, -2),
				CompletedAt: &completedAt,
				Score: &entity.SecurityScore{
					Overall:    80.0,
					Blocked:    4,
					Detected:   2,
					Successful: 4,
					Total:      10,
				},
			},
		},
	}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response application.ExecutionSummary
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.TotalExecutions != 1 {
		t.Errorf("TotalExecutions = %d, want 1", response.TotalExecutions)
	}
}

func TestAnalyticsHandler_GetPeriodStats_Success(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	now := time.Now()
	start := url.QueryEscape(now.AddDate(0, 0, -7).Format(time.RFC3339))
	end := url.QueryEscape(now.Format(time.RFC3339))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+start+"&end="+end, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAnalyticsHandler_GetPeriodStats_MissingStart(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?end="+url.QueryEscape(time.Now().Format(time.RFC3339)), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetPeriodStats_MissingEnd(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+url.QueryEscape(time.Now().Format(time.RFC3339)), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetPeriodStats_InvalidStartFormat(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start=invalid&end="+url.QueryEscape(time.Now().Format(time.RFC3339)), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetPeriodStats_InvalidEndFormat(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+url.QueryEscape(time.Now().Format(time.RFC3339))+"&end=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetPeriodStats_EndBeforeStart(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	now := time.Now()
	end := url.QueryEscape(now.AddDate(0, 0, -7).Format(time.RFC3339))
	start := url.QueryEscape(now.Format(time.RFC3339))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+start+"&end="+end, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAnalyticsHandler_CompareScores_InvalidDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	// Invalid days should fall back to default
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_CompareScores_DaysOutOfRange(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	// Days > 365 should fall back to default
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=1000", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// mockResultRepoForHandler implements repository.ResultRepository for handler tests
type mockResultRepoForHandler struct {
	executions []*entity.Execution
}

func (m *mockResultRepoForHandler) CreateExecution(ctx context.Context, execution *entity.Execution) error {
	return nil
}

func (m *mockResultRepoForHandler) UpdateExecution(ctx context.Context, execution *entity.Execution) error {
	return nil
}

func (m *mockResultRepoForHandler) FindExecutionByID(ctx context.Context, id string) (*entity.Execution, error) {
	return nil, nil
}

func (m *mockResultRepoForHandler) FindExecutionsByScenario(ctx context.Context, scenarioID string) ([]*entity.Execution, error) {
	return nil, nil
}

func (m *mockResultRepoForHandler) FindRecentExecutions(ctx context.Context, limit int) ([]*entity.Execution, error) {
	return m.executions, nil
}

func (m *mockResultRepoForHandler) FindExecutionsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Execution, error) {
	var results []*entity.Execution
	for _, e := range m.executions {
		if !e.StartedAt.Before(start) && !e.StartedAt.After(end) {
			results = append(results, e)
		}
	}
	return results, nil
}

func (m *mockResultRepoForHandler) FindCompletedExecutionsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Execution, error) {
	var results []*entity.Execution
	for _, e := range m.executions {
		if e.Status == entity.ExecutionCompleted && !e.StartedAt.Before(start) && !e.StartedAt.After(end) {
			results = append(results, e)
		}
	}
	return results, nil
}

func (m *mockResultRepoForHandler) CreateResult(ctx context.Context, result *entity.ExecutionResult) error {
	return nil
}

func (m *mockResultRepoForHandler) UpdateResult(ctx context.Context, result *entity.ExecutionResult) error {
	return nil
}

func (m *mockResultRepoForHandler) FindResultByID(ctx context.Context, id string) (*entity.ExecutionResult, error) {
	return nil, nil
}

func (m *mockResultRepoForHandler) FindResultsByExecution(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	return nil, nil
}

func (m *mockResultRepoForHandler) FindResultsByTechnique(ctx context.Context, techniqueID string) ([]*entity.ExecutionResult, error) {
	return nil, nil
}

// Tests for unauthenticated access
func TestAnalyticsHandler_CompareScores_Unauthenticated(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/compare", handler.CompareScores)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_Unauthenticated(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/trend", handler.GetScoreTrend)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetExecutionSummary_Unauthenticated(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/summary", handler.GetExecutionSummary)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetPeriodStats_Unauthenticated(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/period", handler.GetPeriodStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start=2025-01-01T00:00:00Z&end=2025-01-02T00:00:00Z", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_InvalidDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_DaysOutOfRange(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=1000", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_NegativeDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=-5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetExecutionSummary_WithDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary?days=7", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAnalyticsHandler_GetExecutionSummary_InvalidDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary?days=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_CompareScores_NegativeDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=-10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetExecutionSummary_DaysOutOfRange(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	service := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(service)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary?days=500", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// mockErrorResultRepoForHandler implements repository.ResultRepository and always returns errors
type mockErrorResultRepoForHandler struct {
	err error
}

func (m *mockErrorResultRepoForHandler) CreateExecution(ctx context.Context, execution *entity.Execution) error {
	return m.err
}
func (m *mockErrorResultRepoForHandler) UpdateExecution(ctx context.Context, execution *entity.Execution) error {
	return m.err
}
func (m *mockErrorResultRepoForHandler) FindExecutionByID(ctx context.Context, id string) (*entity.Execution, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) FindExecutionsByScenario(ctx context.Context, scenarioID string) ([]*entity.Execution, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) FindRecentExecutions(ctx context.Context, limit int) ([]*entity.Execution, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) FindExecutionsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Execution, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) FindCompletedExecutionsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Execution, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) CreateResult(ctx context.Context, result *entity.ExecutionResult) error {
	return m.err
}
func (m *mockErrorResultRepoForHandler) UpdateResult(ctx context.Context, result *entity.ExecutionResult) error {
	return m.err
}
func (m *mockErrorResultRepoForHandler) FindResultByID(ctx context.Context, id string) (*entity.ExecutionResult, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) FindResultsByExecution(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	return nil, m.err
}
func (m *mockErrorResultRepoForHandler) FindResultsByTechnique(ctx context.Context, techniqueID string) ([]*entity.ExecutionResult, error) {
	return nil, m.err
}

// --- Service error path tests ---

func TestAnalyticsHandler_CompareScores_ServiceError(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database connection failed")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to compare scores" {
		t.Errorf("Expected error 'failed to compare scores', got %q", response["error"])
	}
}

func TestAnalyticsHandler_GetScoreTrend_ServiceError(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database connection failed")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get score trend" {
		t.Errorf("Expected error 'failed to get score trend', got %q", response["error"])
	}
}

func TestAnalyticsHandler_GetExecutionSummary_ServiceError(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database connection failed")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get execution summary" {
		t.Errorf("Expected error 'failed to get execution summary', got %q", response["error"])
	}
}

func TestAnalyticsHandler_GetPeriodStats_ServiceError(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database connection failed")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	now := time.Now()
	start := url.QueryEscape(now.AddDate(0, 0, -7).Format(time.RFC3339))
	end := url.QueryEscape(now.Format(time.RFC3339))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+start+"&end="+end, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get period stats" {
		t.Errorf("Expected error 'failed to get period stats', got %q", response["error"])
	}
}

// --- Invalid query param edge cases ---

func TestAnalyticsHandler_GetPeriodStats_BothDatesMissing(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "start and end dates are required" {
		t.Errorf("Expected error 'start and end dates are required', got %q", response["error"])
	}
}

func TestAnalyticsHandler_CompareScores_DaysZero(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	// days=0 should fall back to default (d > 0 check fails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=0", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_DaysZero(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	// days=0 should fall back to default (d > 0 check fails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=0", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetExecutionSummary_DaysZero(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	// days=0 should fall back to default (d > 0 check fails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary?days=0", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_CompareScores_DaysExactly365(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	// days=365 is exactly at the boundary (d <= 365), should be accepted
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=365", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_DaysExactly365(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	// days=365 is exactly at the boundary (d <= 365), should be accepted
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=365", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response application.ScoreTrend
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	// 365 doesn't match 30 or 90, so it falls through to "7d" in periodLabel
	if response.Period != "7d" {
		t.Errorf("Period = %q, want '7d'", response.Period)
	}
}

func TestAnalyticsHandler_GetExecutionSummary_NegativeDays(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	// Negative days should fall back to default (d > 0 check fails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary?days=-10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetPeriodStats_StartDateNotRFC3339(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	// Use a date format that is valid but not RFC3339 (e.g., "2025-01-01")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start=2025-01-01&end="+url.QueryEscape(time.Now().Format(time.RFC3339)), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "invalid start date format" {
		t.Errorf("Expected error 'invalid start date format', got %q", response["error"])
	}
}

func TestAnalyticsHandler_GetPeriodStats_EndDateNotRFC3339(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	// Valid start, but end is not RFC3339
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+url.QueryEscape(time.Now().Format(time.RFC3339))+"&end=2025-12-31", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "invalid end date format" {
		t.Errorf("Expected error 'invalid end date format', got %q", response["error"])
	}
}

func TestAnalyticsHandler_GetPeriodStats_SameDates(t *testing.T) {
	repo := &mockResultRepoForHandler{}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/period", withAuthAnalytics(handler.GetPeriodStats))

	// Same start and end date - end is not before start so should succeed
	now := time.Now()
	dateStr := url.QueryEscape(now.Format(time.RFC3339))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/period?start="+dateStr+"&end="+dateStr, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAnalyticsHandler_CompareScores_ServiceErrorWithCustomDays(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database timeout")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/compare", withAuthAnalytics(handler.CompareScores))

	// Service error even with valid custom days
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/compare?days=14", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetScoreTrend_ServiceErrorWithCustomDays(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database timeout")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/trend", withAuthAnalytics(handler.GetScoreTrend))

	// Service error even with valid custom days
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trend?days=14", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetExecutionSummary_ServiceErrorWithCustomDays(t *testing.T) {
	repo := &mockErrorResultRepoForHandler{err: errors.New("database timeout")}
	svc := application.NewAnalyticsService(repo)
	handler := NewAnalyticsHandler(svc)

	router := gin.New()
	router.GET("/summary", withAuthAnalytics(handler.GetExecutionSummary))

	// Service error even with valid custom days
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/summary?days=14", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
