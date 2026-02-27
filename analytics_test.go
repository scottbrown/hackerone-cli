package hackeronecli

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetAnalytics(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/analytics" {
			t.Errorf("expected path /analytics, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("program") != "my-program" {
			t.Errorf("expected program=my-program, got %q", r.URL.Query().Get("program"))
		}
		if r.URL.Query().Get("groups") != "severity" {
			t.Errorf("expected groups=severity, got %q", r.URL.Query().Get("groups"))
		}
		if r.URL.Query().Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %q", r.URL.Query().Get("start_date"))
		}
		if r.URL.Query().Get("end_date") != "2024-12-31" {
			t.Errorf("expected end_date=2024-12-31, got %q", r.URL.Query().Get("end_date"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"total": "42"},
		})
	})

	result, err := c.GetAnalytics(context.Background(), "my-program", "severity", "2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestGetAnalyticsMinimalParams(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("program") != "test-prog" {
			t.Errorf("expected program=test-prog, got %q", r.URL.Query().Get("program"))
		}
		if r.URL.Query().Get("groups") != "" {
			t.Errorf("expected no groups param, got %q", r.URL.Query().Get("groups"))
		}
		if r.URL.Query().Get("start_date") != "" {
			t.Errorf("expected no start_date param, got %q", r.URL.Query().Get("start_date"))
		}
		if r.URL.Query().Get("end_date") != "" {
			t.Errorf("expected no end_date param, got %q", r.URL.Query().Get("end_date"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []string{},
		})
	})

	result, err := c.GetAnalytics(context.Background(), "test-prog", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGetAnalyticsError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"})
	})

	_, err := c.GetAnalytics(context.Background(), "my-program", "", "", "")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("expected status 403, got %d", apiErr.StatusCode)
	}
}
