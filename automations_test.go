package hackeronecli

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListAutomations(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/automations" {
			t.Errorf("expected path /automations, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page[number]") != "1" {
			t.Errorf("expected page[number]=1, got %s", r.URL.Query().Get("page[number]"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []Automation{
				{ID: "auto1", Type: "automation", Attributes: AutomationAttributes{Name: "test-auto", State: "active"}},
			},
		})
	})

	automations, err := c.ListAutomations(context.Background(), PageParams{Number: 1, Size: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(automations) != 1 {
		t.Fatalf("expected 1 automation, got %d", len(automations))
	}
	if automations[0].ID != "auto1" {
		t.Errorf("expected id auto1, got %s", automations[0].ID)
	}
	if automations[0].Attributes.Name != "test-auto" {
		t.Errorf("expected name test-auto, got %s", automations[0].Attributes.Name)
	}
}

func TestGetAutomation(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/automations/auto1" {
			t.Errorf("expected path /automations/auto1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": Automation{ID: "auto1", Type: "automation", Attributes: AutomationAttributes{Name: "test-auto"}},
		})
	})

	auto, err := c.GetAutomation(context.Background(), "auto1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auto.ID != "auto1" {
		t.Errorf("expected id auto1, got %s", auto.ID)
	}
}

func TestTriggerAutomation(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/automations/auto1/trigger" {
			t.Errorf("expected path /automations/auto1/trigger, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	err := c.TriggerAutomation(context.Background(), "auto1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListAutomationRuns(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/automations/auto1/runs" {
			t.Errorf("expected path /automations/auto1/runs, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []AutomationRun{
				{ID: "run1", Type: "automation_run", Attributes: AutomationRunAttributes{Status: "completed"}},
			},
		})
	})

	runs, err := c.ListAutomationRuns(context.Background(), "auto1", PageParams{Number: 1, Size: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].Attributes.Status != "completed" {
		t.Errorf("expected status completed, got %s", runs[0].Attributes.Status)
	}
}

func TestGetAutomationRun(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/automations/auto1/runs/run1" {
			t.Errorf("expected path /automations/auto1/runs/run1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": AutomationRun{ID: "run1", Type: "automation_run", Attributes: AutomationRunAttributes{Status: "completed"}},
		})
	})

	run, err := c.GetAutomationRun(context.Background(), "auto1", "run1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if run.ID != "run1" {
		t.Errorf("expected id run1, got %s", run.ID)
	}
}

func TestGetAutomationRunLogs(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/automations/auto1/runs/run1/logs" {
			t.Errorf("expected path /automations/auto1/runs/run1/logs, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []AutomationRunLog{
				{ID: "log1", Type: "automation_run_log", Attributes: AutomationRunLogAttributes{Message: "started", Level: "info"}},
			},
		})
	})

	logs, err := c.GetAutomationRunLogs(context.Background(), "auto1", "run1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].Attributes.Message != "started" {
		t.Errorf("expected message 'started', got %s", logs[0].Attributes.Message)
	}
}

func TestListAutomationsAPIError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"})
	})

	_, err := c.ListAutomations(context.Background(), PageParams{})
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

func TestTriggerAutomationAPIError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	})

	err := c.TriggerAutomation(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAutomationRunLogsAPIError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
	})

	_, err := c.GetAutomationRunLogs(context.Background(), "auto1", "run1")
	if err == nil {
		t.Fatal("expected error")
	}
}
