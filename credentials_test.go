package hackeronecli

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListCredentials(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/credentials" {
			t.Errorf("expected path /credentials, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page[number]") != "1" {
			t.Errorf("expected page[number]=1, got %q", r.URL.Query().Get("page[number]"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "cred-1", "type": "credential", "attributes": map[string]string{"account_name": "acct1", "state": "active"}},
			},
		})
	})

	creds, err := c.ListCredentials(context.Background(), PageParams{Number: 1, Size: 25})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(creds) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(creds))
	}
	if creds[0].ID != "cred-1" {
		t.Errorf("expected id cred-1, got %q", creds[0].ID)
	}
	if creds[0].Attributes.AccountName != "acct1" {
		t.Errorf("expected account_name acct1, got %q", creds[0].Attributes.AccountName)
	}
}

func TestCreateCredential(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/credentials" {
			t.Errorf("expected path /credentials, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var input CreateCredentialInput
		json.Unmarshal(body, &input)
		if input.AccountName != "acct1" {
			t.Errorf("expected account_name acct1, got %q", input.AccountName)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "cred-new", "type": "credential", "attributes": map[string]string{"account_name": "acct1"}},
		})
	})

	cred, err := c.CreateCredential(context.Background(), CreateCredentialInput{
		AccountName: "acct1",
		CredType:    "password",
		Credentials: "secret",
		ProgramID:   "prog-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.ID != "cred-new" {
		t.Errorf("expected id cred-new, got %q", cred.ID)
	}
}

func TestUpdateCredential(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/credentials/cred-1" {
			t.Errorf("expected path /credentials/cred-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "cred-1", "type": "credential", "attributes": map[string]string{"account_name": "updated"}},
		})
	})

	cred, err := c.UpdateCredential(context.Background(), "cred-1", UpdateCredentialInput{AccountName: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.Attributes.AccountName != "updated" {
		t.Errorf("expected account_name updated, got %q", cred.Attributes.AccountName)
	}
}

func TestDeleteCredential(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/credentials/cred-1" {
			t.Errorf("expected path /credentials/cred-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteCredential(context.Background(), "cred-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssignCredential(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/credentials/cred-1/assign" {
			t.Errorf("expected path /credentials/cred-1/assign, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "cred-1", "type": "credential", "attributes": map[string]string{"state": "assigned"}},
		})
	})

	cred, err := c.AssignCredential(context.Background(), "cred-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.Attributes.State != "assigned" {
		t.Errorf("expected state assigned, got %q", cred.Attributes.State)
	}
}

func TestRevokeCredential(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/credentials/cred-1/revoke" {
			t.Errorf("expected path /credentials/cred-1/revoke, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "cred-1", "type": "credential", "attributes": map[string]string{"state": "revoked"}},
		})
	})

	cred, err := c.RevokeCredential(context.Background(), "cred-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.Attributes.State != "revoked" {
		t.Errorf("expected state revoked, got %q", cred.Attributes.State)
	}
}

func TestListCredentialInquiries(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/prog-1/credential_inquiries" {
			t.Errorf("expected path /programs/prog-1/credential_inquiries, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "inq-1", "type": "credential_inquiry", "attributes": map[string]interface{}{"question": "API key?", "required": true}},
			},
		})
	})

	inquiries, err := c.ListCredentialInquiries(context.Background(), "prog-1", PageParams{Number: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inquiries) != 1 {
		t.Fatalf("expected 1 inquiry, got %d", len(inquiries))
	}
	if inquiries[0].ID != "inq-1" {
		t.Errorf("expected id inq-1, got %q", inquiries[0].ID)
	}
}

func TestCreateCredentialInquiry(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/programs/prog-1/credential_inquiries" {
			t.Errorf("expected path /programs/prog-1/credential_inquiries, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "inq-new", "type": "credential_inquiry", "attributes": map[string]interface{}{"question": "Token?"}},
		})
	})

	inq, err := c.CreateCredentialInquiry(context.Background(), "prog-1", CreateCredentialInquiryInput{
		Question:  "Token?",
		Required:  true,
		FieldType: "text",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inq.ID != "inq-new" {
		t.Errorf("expected id inq-new, got %q", inq.ID)
	}
}

func TestUpdateCredentialInquiry(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/programs/prog-1/credential_inquiries/inq-1" {
			t.Errorf("expected path /programs/prog-1/credential_inquiries/inq-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "inq-1", "type": "credential_inquiry", "attributes": map[string]interface{}{"question": "Updated?"}},
		})
	})

	inq, err := c.UpdateCredentialInquiry(context.Background(), "prog-1", "inq-1", CreateCredentialInquiryInput{
		Question:  "Updated?",
		Required:  false,
		FieldType: "text",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inq.Attributes.Question != "Updated?" {
		t.Errorf("expected question Updated?, got %q", inq.Attributes.Question)
	}
}

func TestDeleteCredentialInquiry(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/programs/prog-1/credential_inquiries/inq-1" {
			t.Errorf("expected path /programs/prog-1/credential_inquiries/inq-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteCredentialInquiry(context.Background(), "prog-1", "inq-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListCredentialInquiryResponses(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/programs/prog-1/credential_inquiries/inq-1/credential_inquiry_responses" {
			t.Errorf("expected correct path, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "resp-1", "type": "credential_inquiry_response"},
			},
		})
	})

	responses, err := c.ListCredentialInquiryResponses(context.Background(), "prog-1", "inq-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}
	if responses[0].ID != "resp-1" {
		t.Errorf("expected id resp-1, got %q", responses[0].ID)
	}
}

func TestDeleteCredentialInquiryResponse(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/programs/prog-1/credential_inquiries/inq-1/credential_inquiry_responses/resp-1" {
			t.Errorf("expected correct path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteCredentialInquiryResponse(context.Background(), "prog-1", "inq-1", "resp-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListCredentialsError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"})
	})

	_, err := c.ListCredentials(context.Background(), PageParams{})
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
