//go:build unit

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/openpayment/gateway/internal/config"
)

func testRouter() *chi.Mux {
	cfg := config.Load()
	return NewRouter(cfg, &HealthChecker{
		Uptime:  time.Now(),
		Version: "test",
	})
}

func TestHealthEndpoint(t *testing.T) {
	router := testRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err)

	data, ok := body["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "ok", data["status"])
	_, hasChecks := data["checks"]
	assert.True(t, hasChecks)
}

func TestHealthEndpoint_MethodNotAllowed(t *testing.T) {
	router := testRouter()

	methods := []string{"POST", "PUT", "PATCH", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}

func TestHealthEndpoint_ContentType(t *testing.T) {
	router := testRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestHealthEndpoint_NotFound(t *testing.T) {
	router := testRouter()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCorsHeaders(t *testing.T) {
	router := testRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Origin"), "")
}

func TestSecurityHeaders(t *testing.T) {
	router := testRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
}
