//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func waitForHealthyStatus(t *testing.T, url string, timeout time.Duration) map[string]any {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp != nil {
			if resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()
				var body map[string]any
				if err := json.NewDecoder(resp.Body).Decode(&body); err == nil {
					return body
				}
			}
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("service did not become healthy in time: %s", url)
	return nil
}

func TestAPIStatusEndpoints(t *testing.T) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	t.Run("v1 status", func(t *testing.T) {
		body := waitForHealthyStatus(t, fmt.Sprintf("%s/api/v1/status", baseURL), 20*time.Second)
		if body["status"] != "ok" {
			t.Fatalf("expected status ok, got %v", body["status"])
		}
	})

	t.Run("v2 status", func(t *testing.T) {
		body := waitForHealthyStatus(t, fmt.Sprintf("%s/api/v2/status", baseURL), 20*time.Second)
		if body["status"] != "ok" {
			t.Fatalf("expected status ok, got %v", body["status"])
		}
		if body["version"] != "v2" {
			t.Fatalf("expected version v2, got %v", body["version"])
		}
	})
}
