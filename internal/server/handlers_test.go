package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func doGet(t *testing.T, handler gin.HandlerFunc, path string) *httptest.ResponseRecorder {
	t.Helper()
	r := gin.New()
	r.GET(path, handler)
	req := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestWeatherHandler(t *testing.T) {
	w := doGet(t, WeatherHandler, "/weather")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("invalid JSON response")
	}

	requiredFields := []string{"city", "temperature", "condition", "humidity", "timestamp"}
	for _, field := range requiredFields {
		if _, ok := resp[field]; !ok {
			t.Errorf("missing field: %s", field)
		}
	}
}

func TestJokeHandler(t *testing.T) {
	w := doGet(t, JokeHandler, "/joke")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("invalid JSON response")
	}

	if _, ok := resp["setup"]; !ok {
		t.Error("missing field: setup")
	}
	if _, ok := resp["punchline"]; !ok {
		t.Error("missing field: punchline")
	}
}

func TestPremiumDataHandler(t *testing.T) {
	w := doGet(t, PremiumDataHandler, "/premium-data")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("invalid JSON response")
	}

	if _, ok := resp["report"]; !ok {
		t.Error("missing field: report")
	}
	if _, ok := resp["metrics"]; !ok {
		t.Error("missing field: metrics")
	}
	if _, ok := resp["generatedAt"]; !ok {
		t.Error("missing field: generatedAt")
	}
}

func TestHealthHandler(t *testing.T) {
	handler := HealthHandler("test-service", "eip155:84532")
	w := doGet(t, handler, "/health")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", resp["status"])
	}
	if resp["service"] != "test-service" {
		t.Errorf("expected service=test-service, got %v", resp["service"])
	}
	if resp["network"] != "eip155:84532" {
		t.Errorf("expected network=eip155:84532, got %v", resp["network"])
	}
}
