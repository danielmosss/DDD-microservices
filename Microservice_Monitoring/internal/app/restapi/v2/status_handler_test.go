package v2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetStatus_ReturnsContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v2/status", GetStatus)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/status", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json body: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status field 'ok', got '%s'", body["status"])
	}
	if body["version"] != "v2" {
		t.Fatalf("expected version field 'v2', got '%s'", body["version"])
	}
	if body["message"] == "" {
		t.Fatal("expected non-empty message field")
	}
}
