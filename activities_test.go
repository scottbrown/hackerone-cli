package hackeronecli

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetActivity(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/activities/abc123" {
			t.Errorf("expected path /activities/abc123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(activityResponse{
			Data: Activity{
				ID:   "abc123",
				Type: "activity-comment",
				Attributes: ActivityAttributes{
					Message:   "test message",
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-02T00:00:00Z",
				},
			},
		})
	})

	activity, err := c.GetActivity(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if activity.ID != "abc123" {
		t.Errorf("expected ID abc123, got %q", activity.ID)
	}
	if activity.Type != "activity-comment" {
		t.Errorf("expected type activity-comment, got %q", activity.Type)
	}
	if activity.Attributes.Message != "test message" {
		t.Errorf("expected message 'test message', got %q", activity.Attributes.Message)
	}
}

func TestGetActivityError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	})

	_, err := c.GetActivity(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

func TestListActivities(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/incremental/activities" {
			t.Errorf("expected path /incremental/activities, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page[number]") != "1" {
			t.Errorf("expected page[number]=1, got %q", r.URL.Query().Get("page[number]"))
		}
		if r.URL.Query().Get("page[size]") != "25" {
			t.Errorf("expected page[size]=25, got %q", r.URL.Query().Get("page[size]"))
		}
		if r.URL.Query().Get("page[updated_at_after]") != "2024-01-01" {
			t.Errorf("expected updated_at_after=2024-01-01, got %q", r.URL.Query().Get("page[updated_at_after]"))
		}
		if r.URL.Query().Get("page[updated_at_before]") != "2024-12-31" {
			t.Errorf("expected updated_at_before=2024-12-31, got %q", r.URL.Query().Get("page[updated_at_before]"))
		}
		json.NewEncoder(w).Encode(activitiesResponse{
			Data: []Activity{
				{ID: "1", Type: "activity-comment"},
				{ID: "2", Type: "activity-bug-filed"},
			},
		})
	})

	params := PageParams{Number: 1, Size: 25}
	activities, err := c.ListActivities(context.Background(), params, "2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(activities) != 2 {
		t.Fatalf("expected 2 activities, got %d", len(activities))
	}
	if activities[0].ID != "1" {
		t.Errorf("expected first activity ID 1, got %q", activities[0].ID)
	}
	if activities[1].ID != "2" {
		t.Errorf("expected second activity ID 2, got %q", activities[1].ID)
	}
}

func TestListActivitiesNoFilters(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page[updated_at_after]") != "" {
			t.Errorf("expected no updated_at_after param, got %q", r.URL.Query().Get("page[updated_at_after]"))
		}
		if r.URL.Query().Get("page[updated_at_before]") != "" {
			t.Errorf("expected no updated_at_before param, got %q", r.URL.Query().Get("page[updated_at_before]"))
		}
		json.NewEncoder(w).Encode(activitiesResponse{Data: []Activity{}})
	})

	activities, err := c.ListActivities(context.Background(), PageParams{}, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(activities) != 0 {
		t.Errorf("expected 0 activities, got %d", len(activities))
	}
}

func TestListActivitiesError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
	})

	_, err := c.ListActivities(context.Background(), PageParams{}, "", "")
	if err == nil {
		t.Fatal("expected error")
	}
}
