package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fresh132/URL-check/internal/api"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/check", api.Check)
	r.POST("/report", api.Report)
	return r
}

func TestCheck_ValidRequest(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string][]string{"url": {"https://example.com"}})
	req, _ := http.NewRequest("POST", "/check", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"id"`)) {
		t.Errorf("response missing id field")
	}
}

func TestCheck_EmptyRequest(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string][]string{"url": {}})
	req, _ := http.NewRequest("POST", "/check", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestReport_InvalidRequest(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string][]string{"links_num": {}})
	req, _ := http.NewRequest("POST", "/report", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
