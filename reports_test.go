package hackeronecli

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListReports(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reports" {
			t.Errorf("expected path /reports, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page[number]") != "1" {
			t.Errorf("expected page[number]=1, got %s", r.URL.Query().Get("page[number]"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "1", "type": "report", "attributes": map[string]string{"title": "XSS Bug"}},
			},
		})
	})

	reports, err := c.ListReports(context.Background(), PageParams{Number: 1, Size: 25}, ListReportsFilter{Programs: []string{"my-program"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}
	if reports[0].ID != "1" {
		t.Errorf("expected report ID '1', got %q", reports[0].ID)
	}
}

func TestGetReport(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123" {
			t.Errorf("expected path /reports/123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "123", "type": "report", "attributes": map[string]string{"title": "Test"}},
		})
	})

	report, err := c.GetReport(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ID != "123" {
		t.Errorf("expected ID '123', got %q", report.ID)
	}
}

func TestCreateReport(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports" {
			t.Errorf("expected path /reports, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var input CreateReportInput
		json.Unmarshal(body, &input)
		if input.Title != "XSS" {
			t.Errorf("expected title 'XSS', got %q", input.Title)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "456", "type": "report"},
		})
	})

	report, err := c.CreateReport(context.Background(), CreateReportInput{
		Title:     "XSS",
		ProgramID: "prog-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.ID != "456" {
		t.Errorf("expected ID '456', got %q", report.ID)
	}
}

func TestAddComment(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/comments" {
			t.Errorf("expected path /reports/123/comments, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var input CommentInput
		json.Unmarshal(body, &input)
		if input.Message != "test comment" {
			t.Errorf("expected message 'test comment', got %q", input.Message)
		}
		if !input.Internal {
			t.Error("expected internal=true")
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	result, err := c.AddComment(context.Background(), "123", CommentInput{Message: "test comment", Internal: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", result["status"])
	}
}

func TestUpdateAssignee(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/assignee" {
			t.Errorf("expected path /reports/123/assignee, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateAssignee(context.Background(), "123", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloseComments(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/comments/close" {
			t.Errorf("expected path /reports/123/comments/close, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.CloseComments(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCustomFields(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/custom_fields" {
			t.Errorf("expected path /reports/123/custom_fields, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateCustomFields(context.Background(), "123", map[string]interface{}{"field1": "value1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCVEs(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/cves" {
			t.Errorf("expected path /reports/123/cves, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateCVEs(context.Background(), "123", []string{"CVE-2024-0001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateInboxes(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/inboxes" {
			t.Errorf("expected path /reports/123/inboxes, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateInboxes(context.Background(), "123", []string{"inbox-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateSeverity(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/severity" {
			t.Errorf("expected path /reports/123/severity, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateSeverity(context.Background(), "123", SeverityInput{Rating: "high", Score: 8.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangeState(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/state" {
			t.Errorf("expected path /reports/123/state, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var input StateChangeInput
		json.Unmarshal(body, &input)
		if input.State != "triaged" {
			t.Errorf("expected state 'triaged', got %q", input.State)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.ChangeState(context.Background(), "123", StateChangeInput{State: "triaged", Message: "confirmed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateReportScope(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/structured_scope" {
			t.Errorf("expected path /reports/123/structured_scope, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateReportScope(context.Background(), "123", "scope-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTitle(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/title" {
			t.Errorf("expected path /reports/123/title, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateTitle(context.Background(), "123", "New Title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateWeakness(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/weakness" {
			t.Errorf("expected path /reports/123/weakness, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateWeakness(context.Background(), "123", "weakness-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateReference(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/reference" {
			t.Errorf("expected path /reports/123/reference, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateReference(context.Background(), "123", "JIRA-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRedactReport(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/redact" {
			t.Errorf("expected path /reports/123/redact, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.RedactReport(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddSummary(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/summary" {
			t.Errorf("expected path /reports/123/summary, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.AddSummary(context.Background(), "123", "Summary content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGeneratePDF(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/pdf" {
			t.Errorf("expected path /reports/123/pdf, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"url": "https://example.com/report.pdf"})
	})

	result, err := c.GeneratePDF(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["url"] != "https://example.com/report.pdf" {
		t.Errorf("expected url, got %v", result["url"])
	}
}

func TestTransferReport(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/transfer" {
			t.Errorf("expected path /reports/123/transfer, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.TransferReport(context.Background(), "123", "prog-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEscalateReport(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/escalate" {
			t.Errorf("expected path /reports/123/escalate, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.EscalateReport(context.Background(), "123", "jira")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeescalateReport(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/escalations" {
			t.Errorf("expected path /reports/123/escalations, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeescalateReport(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddParticipant(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/participants" {
			t.Errorf("expected path /reports/123/participants, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.AddParticipant(context.Background(), "123", "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUploadAttachment(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("test content"), 0644)

	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/attachments" {
			t.Errorf("expected path /reports/123/attachments, got %s", r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart content type, got %s", ct)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "att-1"})
	})

	result, err := c.UploadAttachment(context.Background(), "123", tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "att-1" {
		t.Errorf("expected id 'att-1', got %v", result["id"])
	}
}

func TestUploadAttachmentFileNotFound(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	})

	_, err := c.UploadAttachment(context.Background(), "123", "/nonexistent/file.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDeleteAttachment(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/attachments/att-1" {
			t.Errorf("expected path /reports/123/attachments/att-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteAttachment(context.Background(), "123", "att-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAwardReportBounty(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/bounties" {
			t.Errorf("expected path /reports/123/bounties, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.AwardReportBounty(context.Background(), "123", BountyInput{Amount: 500, BonusAmount: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarkIneligible(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/eligibility" {
			t.Errorf("expected path /reports/123/eligibility, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.MarkIneligible(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListBountySuggestions(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/bounty_suggestions" {
			t.Errorf("expected path /reports/123/bounty_suggestions, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "bs-1", "type": "bounty-suggestion"},
			},
		})
	})

	suggestions, err := c.ListBountySuggestions(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestCreateBountySuggestion(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/bounty_suggestions" {
			t.Errorf("expected path /reports/123/bounty_suggestions, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "bs-2", "type": "bounty-suggestion"},
		})
	})

	suggestion, err := c.CreateBountySuggestion(context.Background(), "123", CreateBountySuggestionInput{Amount: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suggestion.ID != "bs-2" {
		t.Errorf("expected ID 'bs-2', got %q", suggestion.ID)
	}
}

func TestUpdateDisclosure(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/disclosure" {
			t.Errorf("expected path /reports/123/disclosure, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateDisclosure(context.Background(), "123", "full")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelDisclosure(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/disclosure" {
			t.Errorf("expected path /reports/123/disclosure, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.CancelDisclosure(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTags(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/tags" {
			t.Errorf("expected path /reports/123/tags, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdateTags(context.Background(), "123", []string{"xss", "critical"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequestRetest(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/retests" {
			t.Errorf("expected path /reports/123/retests, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.RequestRetest(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelRetest(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/retests/rt-1" {
			t.Errorf("expected path /reports/123/retests/rt-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.CancelRetest(context.Background(), "123", "rt-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAwardSwag(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reports/123/swag" {
			t.Errorf("expected path /reports/123/swag, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.AwardSwag(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListReportsAPIError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"})
	})

	_, err := c.ListReports(context.Background(), PageParams{}, ListReportsFilter{})
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
