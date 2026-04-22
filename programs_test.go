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

func TestListPrograms(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs" {
			t.Errorf("expected path /programs, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "1", "type": "program", "attributes": map[string]string{"handle": "test-prog"}},
			},
		})
	})

	programs, err := c.ListPrograms(context.Background(), PageParams{Number: 1, Size: 25})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(programs) != 1 {
		t.Fatalf("expected 1 program, got %d", len(programs))
	}
	if programs[0].ID != "1" {
		t.Errorf("expected program ID '1', got %q", programs[0].ID)
	}
}

func TestGetProgram(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123" {
			t.Errorf("expected path /programs/123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "123", "type": "program", "attributes": map[string]string{"handle": "test"}},
		})
	})

	program, err := c.GetProgram(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if program.ID != "123" {
		t.Errorf("expected ID '123', got %q", program.ID)
	}
}

func TestGetProgramAuditLog(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/audit_log" {
			t.Errorf("expected path /programs/123/audit_log, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "log-1", "type": "audit-log-entry"},
			},
		})
	})

	entries, err := c.GetProgramAuditLog(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestGetProgramBalance(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/balance" {
			t.Errorf("expected path /programs/123/balance, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"balance": "10000.00", "currency": "USD"})
	})

	balance, err := c.GetProgramBalance(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Balance != "10000.00" {
		t.Errorf("expected balance '10000.00', got %q", balance.Balance)
	}
	if balance.Currency != "USD" {
		t.Errorf("expected currency 'USD', got %q", balance.Currency)
	}
}

func TestListPaymentTransactions(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/payment_transactions" {
			t.Errorf("expected path /programs/123/payment_transactions, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "tx-1", "type": "payment-transaction"},
			},
		})
	})

	txns, err := c.ListPaymentTransactions(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txns) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txns))
	}
}

func TestListCommonResponses(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/common_responses" {
			t.Errorf("expected path /programs/123/common_responses, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "cr-1", "type": "common-response"},
			},
		})
	})

	responses, err := c.ListCommonResponses(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}
}

func TestListReporters(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/reporters" {
			t.Errorf("expected path /programs/123/reporters, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "rep-1", "type": "reporter"},
			},
		})
	})

	reporters, err := c.ListReporters(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reporters) != 1 {
		t.Fatalf("expected 1 reporter, got %d", len(reporters))
	}
}

func TestListTeamMembers(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/members" {
			t.Errorf("expected path /programs/123/members, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "tm-1", "type": "team-member"},
			},
		})
	})

	members, err := c.ListTeamMembers(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
}

func TestListThanks(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/thanks" {
			t.Errorf("expected path /programs/123/thanks, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "th-1", "type": "thanks"},
			},
		})
	})

	thanks, err := c.ListThanks(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(thanks) != 1 {
		t.Fatalf("expected 1 thanks, got %d", len(thanks))
	}
}

func TestListIntegrations(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/integrations" {
			t.Errorf("expected path /programs/123/integrations, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "int-1", "type": "integration"},
			},
		})
	})

	integrations, err := c.ListIntegrations(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(integrations) != 1 {
		t.Fatalf("expected 1 integration, got %d", len(integrations))
	}
}

func TestListTriageReviews(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/triage_reviews" {
			t.Errorf("expected path /programs/123/triage_reviews, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "tr-1", "type": "triage-review"},
			},
		})
	})

	reviews, err := c.ListTriageReviews(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
}

func TestListWeaknesses(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/weaknesses" {
			t.Errorf("expected path /programs/123/weaknesses, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "w-1", "type": "weakness"},
			},
		})
	})

	weaknesses, err := c.ListWeaknesses(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(weaknesses) != 1 {
		t.Fatalf("expected 1 weakness, got %d", len(weaknesses))
	}
}

func TestNotifyExternalPlatform(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/notify_external_platform" {
			t.Errorf("expected path /programs/123/notify_external_platform, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.NotifyExternalPlatform(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendProgramMessage(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/messages" {
			t.Errorf("expected path /programs/123/messages, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var envelope struct {
			Data struct {
				Attributes MessageInput `json:"attributes"`
			} `json:"data"`
		}
		json.Unmarshal(body, &envelope)
		input := envelope.Data.Attributes
		if input.RecipientHandle != "hacker1" {
			t.Errorf("expected recipient 'hacker1', got %q", input.RecipientHandle)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.SendProgramMessage(context.Background(), "123", MessageInput{RecipientHandle: "hacker1", Message: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAwardProgramBounty(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/bounties" {
			t.Errorf("expected path /programs/123/bounties, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.AwardProgramBounty(context.Background(), "123", ProgramBountyInput{ReportID: "r-1", Amount: 500})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListAllowedReporters(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/allowed_reporters" {
			t.Errorf("expected path /programs/123/allowed_reporters, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "ar-1", "type": "reporter"},
			},
		})
	})

	reporters, err := c.ListAllowedReporters(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reporters) != 1 {
		t.Fatalf("expected 1 reporter, got %d", len(reporters))
	}
}

func TestGetAllowedReportersHistory(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/allowed_reporters_history" {
			t.Errorf("expected path /programs/123/allowed_reporters_history, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"action": "added"},
			},
		})
	})

	history, err := c.GetAllowedReportersHistory(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(history))
	}
}

func TestGetAllowedReporterActivities(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/allowed_reporter_activities" {
			t.Errorf("expected path /programs/123/allowed_reporter_activities, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"activity": "login"},
			},
		})
	})

	activities, err := c.GetAllowedReporterActivities(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}
}

func TestGetAllowedReporterUsernameHistory(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/allowed_reporter_username_history" {
			t.Errorf("expected path /programs/123/allowed_reporter_username_history, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"username": "old_name"},
			},
		})
	})

	history, err := c.GetAllowedReporterUsernameHistory(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(history))
	}
}

func TestListCVERequests(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/cve_requests" {
			t.Errorf("expected path /programs/123/cve_requests, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "cve-1", "type": "cve-request"},
			},
		})
	})

	requests, err := c.ListCVERequests(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(requests))
	}
}

func TestCreateCVERequest(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/cve_requests" {
			t.Errorf("expected path /programs/123/cve_requests, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "cve-2", "type": "cve-request"},
		})
	})

	request, err := c.CreateCVERequest(context.Background(), "123", CreateCVERequestInput{ReportID: "r-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if request.ID != "cve-2" {
		t.Errorf("expected ID 'cve-2', got %q", request.ID)
	}
}

func TestListHackerInvitations(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/hacker_invitations" {
			t.Errorf("expected path /programs/123/hacker_invitations, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "inv-1", "type": "hacker-invitation"},
			},
		})
	})

	invitations, err := c.ListHackerInvitations(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invitations) != 1 {
		t.Fatalf("expected 1 invitation, got %d", len(invitations))
	}
}

func TestCreateHackerInvitation(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/hacker_invitations" {
			t.Errorf("expected path /programs/123/hacker_invitations, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "inv-2", "type": "hacker-invitation"},
		})
	})

	invitation, err := c.CreateHackerInvitation(context.Background(), "123", CreateHackerInvitationInput{Email: "test@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invitation.ID != "inv-2" {
		t.Errorf("expected ID 'inv-2', got %q", invitation.ID)
	}
}

func TestUpdatePolicy(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/policy" {
			t.Errorf("expected path /programs/123/policy, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var envelope struct {
			Data struct {
				Attributes PolicyInput `json:"attributes"`
			} `json:"data"`
		}
		json.Unmarshal(body, &envelope)
		input := envelope.Data.Attributes
		if input.Policy != "new policy" {
			t.Errorf("expected policy 'new policy', got %q", input.Policy)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	_, err := c.UpdatePolicy(context.Background(), "123", PolicyInput{Policy: "new policy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAttachToPolicy(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "policy.pdf")
	os.WriteFile(tmpFile, []byte("pdf content"), 0644)

	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/policy_attachments" {
			t.Errorf("expected path /programs/123/policy_attachments, got %s", r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart content type, got %s", ct)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "att-1"})
	})

	result, err := c.AttachToPolicy(context.Background(), "123", tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "att-1" {
		t.Errorf("expected id 'att-1', got %v", result["id"])
	}
}

func TestAttachToPolicyFileNotFound(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	})

	_, err := c.AttachToPolicy(context.Background(), "123", "/nonexistent/file.pdf")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestListScopes(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/structured_scopes" {
			t.Errorf("expected path /programs/123/structured_scopes, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "sc-1", "type": "structured-scope"},
			},
		})
	})

	scopes, err := c.ListScopes(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scopes) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(scopes))
	}
}

func TestCreateScope(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/structured_scopes" {
			t.Errorf("expected path /programs/123/structured_scopes, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "sc-2", "type": "structured-scope"},
		})
	})

	scope, err := c.CreateScope(context.Background(), "123", CreateScopeInput{
		AssetIdentifier:       "*.example.com",
		AssetType:             "URL",
		EligibleForBounty:     true,
		EligibleForSubmission: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scope.ID != "sc-2" {
		t.Errorf("expected ID 'sc-2', got %q", scope.ID)
	}
}

func TestUpdateProgramScope(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/structured_scopes/sc-1" {
			t.Errorf("expected path /programs/123/structured_scopes/sc-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "sc-1", "type": "structured-scope"},
		})
	})

	boolTrue := true
	scope, err := c.UpdateProgramScope(context.Background(), "123", "sc-1", UpdateScopeInput{
		EligibleForBounty: &boolTrue,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scope.ID != "sc-1" {
		t.Errorf("expected ID 'sc-1', got %q", scope.ID)
	}
}

func TestArchiveScope(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/structured_scopes/sc-1/archive" {
			t.Errorf("expected path /programs/123/structured_scopes/sc-1/archive, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	})

	err := c.ArchiveScope(context.Background(), "123", "sc-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListAwardedSwag(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/awarded_swags" {
			t.Errorf("expected path /programs/123/awarded_swags, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "sw-1", "type": "awarded-swag"},
			},
		})
	})

	swag, err := c.ListAwardedSwag(context.Background(), "123", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(swag) != 1 {
		t.Fatalf("expected 1 swag, got %d", len(swag))
	}
}

func TestUpdateAwardedSwag(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/programs/123/awarded_swags/sw-1" {
			t.Errorf("expected path /programs/123/awarded_swags/sw-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "sw-1", "type": "awarded-swag"},
		})
	})

	swag, err := c.UpdateAwardedSwag(context.Background(), "123", "sw-1", UpdateSwagInput{Sent: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if swag.ID != "sw-1" {
		t.Errorf("expected ID 'sw-1', got %q", swag.ID)
	}
}

func TestListProgramsAPIError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"})
	})

	_, err := c.ListPrograms(context.Background(), PageParams{})
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
