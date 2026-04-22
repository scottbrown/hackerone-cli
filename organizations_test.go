package hackeronecli

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListOrganizations(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/me/organizations" {
			t.Errorf("expected path /me/organizations, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "org-1", "type": "organization", "attributes": map[string]string{"handle": "acme", "name": "ACME Corp"}},
			},
		})
	})

	orgs, err := c.ListOrganizations(context.Background(), PageParams{Number: 1, Size: 25})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0].ID != "org-1" {
		t.Errorf("expected id org-1, got %q", orgs[0].ID)
	}
	if orgs[0].Attributes.Handle != "acme" {
		t.Errorf("expected handle acme, got %q", orgs[0].Attributes.Handle)
	}
}

func TestGetOrganizationAuditLog(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/audit_log" {
			t.Errorf("expected path /organizations/org-1/audit_log, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "log-1", "type": "audit_log_entry", "attributes": map[string]string{"action": "user.login", "actor_name": "admin"}},
			},
		})
	})

	entries, err := c.GetOrganizationAuditLog(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Attributes.Action != "user.login" {
		t.Errorf("expected action user.login, got %q", entries[0].Attributes.Action)
	}
}

func TestListOrganizationPrograms(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/programs" {
			t.Errorf("expected path /organizations/org-1/programs, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "prog-1", "type": "program"},
			},
		})
	})

	progs, err := c.ListOrganizationPrograms(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(progs) != 1 {
		t.Fatalf("expected 1 program, got %d", len(progs))
	}
}

func TestListOrganizationInboxes(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/inboxes" {
			t.Errorf("expected path /organizations/org-1/inboxes, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "inbox-1", "type": "inbox"},
			},
		})
	})

	inboxes, err := c.ListOrganizationInboxes(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inboxes) != 1 {
		t.Fatalf("expected 1 inbox, got %d", len(inboxes))
	}
	if inboxes[0].ID != "inbox-1" {
		t.Errorf("expected id inbox-1, got %q", inboxes[0].ID)
	}
}

func TestListEligibilitySettings(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/eligibility_settings" {
			t.Errorf("expected path /organizations/org-1/eligibility_settings, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "es-1", "type": "eligibility_setting"},
			},
		})
	})

	settings, err := c.ListEligibilitySettings(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(settings))
	}
}

func TestGetEligibilitySetting(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/eligibility_settings/es-1" {
			t.Errorf("expected path /organizations/org-1/eligibility_settings/es-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "es-1", "type": "eligibility_setting"},
		})
	})

	setting, err := c.GetEligibilitySetting(context.Background(), "org-1", "es-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if setting.ID != "es-1" {
		t.Errorf("expected id es-1, got %q", setting.ID)
	}
}

func TestListInvitations(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/invitations" {
			t.Errorf("expected path /organizations/org-1/invitations, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "inv-1", "type": "invitation", "attributes": map[string]string{"email": "user@example.com"}},
			},
		})
	})

	invs, err := c.ListInvitations(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invs) != 1 {
		t.Fatalf("expected 1 invitation, got %d", len(invs))
	}
	if invs[0].Attributes.Email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %q", invs[0].Attributes.Email)
	}
}

func TestCreateInvitation(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/invitations" {
			t.Errorf("expected path /organizations/org-1/invitations, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var envelope struct {
			Data struct {
				Attributes CreateInvitationInput `json:"attributes"`
			} `json:"data"`
		}
		json.Unmarshal(body, &envelope)
		input := envelope.Data.Attributes
		if input.Email != "new@example.com" {
			t.Errorf("expected email new@example.com, got %q", input.Email)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "inv-new", "type": "invitation", "attributes": map[string]string{"email": "new@example.com"}},
		})
	})

	inv, err := c.CreateInvitation(context.Background(), "org-1", CreateInvitationInput{Email: "new@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.ID != "inv-new" {
		t.Errorf("expected id inv-new, got %q", inv.ID)
	}
}

func TestListGroups(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/groups" {
			t.Errorf("expected path /organizations/org-1/groups, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "grp-1", "type": "group", "attributes": map[string]interface{}{"name": "admins", "permissions": []string{"read", "write"}}},
			},
		})
	})

	groups, err := c.ListGroups(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Attributes.Name != "admins" {
		t.Errorf("expected name admins, got %q", groups[0].Attributes.Name)
	}
}

func TestGetGroup(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/groups/grp-1" {
			t.Errorf("expected path /organizations/org-1/groups/grp-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "grp-1", "type": "group", "attributes": map[string]interface{}{"name": "admins"}},
		})
	})

	group, err := c.GetGroup(context.Background(), "org-1", "grp-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "grp-1" {
		t.Errorf("expected id grp-1, got %q", group.ID)
	}
}

func TestCreateGroup(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/groups" {
			t.Errorf("expected path /organizations/org-1/groups, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "grp-new", "type": "group", "attributes": map[string]interface{}{"name": "editors"}},
		})
	})

	group, err := c.CreateGroup(context.Background(), "org-1", CreateGroupInput{Name: "editors", Permissions: []string{"read"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "grp-new" {
		t.Errorf("expected id grp-new, got %q", group.ID)
	}
}

func TestUpdateGroup(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/groups/grp-1" {
			t.Errorf("expected path /organizations/org-1/groups/grp-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "grp-1", "type": "group", "attributes": map[string]interface{}{"name": "updated"}},
		})
	})

	group, err := c.UpdateGroup(context.Background(), "org-1", "grp-1", UpdateGroupInput{Name: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Attributes.Name != "updated" {
		t.Errorf("expected name updated, got %q", group.Attributes.Name)
	}
}

func TestListMembers(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/members" {
			t.Errorf("expected path /organizations/org-1/members, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "mem-1", "type": "member", "attributes": map[string]interface{}{"email": "alice@example.com", "name": "Alice"}},
			},
		})
	})

	members, err := c.ListMembers(context.Background(), "org-1", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
	if members[0].Attributes.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %q", members[0].Attributes.Email)
	}
}

func TestGetMember(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/org-1/members/mem-1" {
			t.Errorf("expected path /organizations/org-1/members/mem-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "mem-1", "type": "member", "attributes": map[string]interface{}{"name": "Alice"}},
		})
	})

	member, err := c.GetMember(context.Background(), "org-1", "mem-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.ID != "mem-1" {
		t.Errorf("expected id mem-1, got %q", member.ID)
	}
}

func TestUpdateMember(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/members/mem-1" {
			t.Errorf("expected path /organizations/org-1/members/mem-1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "mem-1", "type": "member", "attributes": map[string]interface{}{"name": "Alice"}},
		})
	})

	member, err := c.UpdateMember(context.Background(), "org-1", "mem-1", UpdateMemberInput{GroupIDs: []string{"grp-1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.ID != "mem-1" {
		t.Errorf("expected id mem-1, got %q", member.ID)
	}
}

func TestRemoveMember(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/members/mem-1" {
			t.Errorf("expected path /organizations/org-1/members/mem-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.RemoveMember(context.Background(), "org-1", "mem-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListOrganizationsError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
	})

	_, err := c.ListOrganizations(context.Background(), PageParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
}
