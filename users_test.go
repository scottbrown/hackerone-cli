package hackeronecli

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/users/johndoe" {
			t.Errorf("expected path /users/johndoe, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(userResponse{
			Data: User{
				ID:       "123",
				Username: "johndoe",
				Name:     "John Doe",
				Bio:      "Security researcher",
				Website:  "https://example.com",
			},
		})
	})

	user, err := c.GetUser(context.Background(), "johndoe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "123" {
		t.Errorf("expected ID 123, got %q", user.ID)
	}
	if user.Username != "johndoe" {
		t.Errorf("expected username johndoe, got %q", user.Username)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %q", user.Name)
	}
	if user.Bio != "Security researcher" {
		t.Errorf("expected bio 'Security researcher', got %q", user.Bio)
	}
	if user.Website != "https://example.com" {
		t.Errorf("expected website https://example.com, got %q", user.Website)
	}
}

func TestGetUserError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	})

	_, err := c.GetUser(context.Background(), "nonexistent")
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

func TestGetUserByID(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/users/456" {
			t.Errorf("expected path /users/456, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(userResponse{
			Data: User{
				ID:       "456",
				Username: "janedoe",
				Name:     "Jane Doe",
			},
		})
	})

	user, err := c.GetUserByID(context.Background(), "456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "456" {
		t.Errorf("expected ID 456, got %q", user.ID)
	}
	if user.Username != "janedoe" {
		t.Errorf("expected username janedoe, got %q", user.Username)
	}
}

func TestGetUserByIDError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal error"})
	})

	_, err := c.GetUserByID(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
}
