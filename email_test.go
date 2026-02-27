package hackeronecli

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestSendEmail(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/email" {
			t.Errorf("expected path /email, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var input SendEmailInput
		if err := json.Unmarshal(body, &input); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}
		if input.To != "user@example.com" {
			t.Errorf("expected to=user@example.com, got %q", input.To)
		}
		if input.Subject != "Test Subject" {
			t.Errorf("expected subject=Test Subject, got %q", input.Subject)
		}
		if input.Body != "Hello, world!" {
			t.Errorf("expected body=Hello, world!, got %q", input.Body)
		}
		w.WriteHeader(http.StatusCreated)
	})

	input := SendEmailInput{
		To:      "user@example.com",
		Subject: "Test Subject",
		Body:    "Hello, world!",
	}
	err := c.SendEmail(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendEmailError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid email"})
	})

	input := SendEmailInput{
		To:      "invalid",
		Subject: "Test",
		Body:    "test",
	}
	err := c.SendEmail(context.Background(), input)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}

func TestSendEmailEmptyBody(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var input SendEmailInput
		if err := json.Unmarshal(body, &input); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}
		if input.Body != "" {
			t.Errorf("expected empty body, got %q", input.Body)
		}
		w.WriteHeader(http.StatusCreated)
	})

	input := SendEmailInput{
		To:      "user@example.com",
		Subject: "No body",
	}
	err := c.SendEmail(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
