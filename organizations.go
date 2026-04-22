package hackeronecli

import (
	"bytes"
	"context"
	"fmt"
)

type Organization struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes OrganizationAttributes `json:"attributes"`
}

type OrganizationAttributes struct {
	Handle    string `json:"handle"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type AuditLogEntry struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes AuditLogAttributes `json:"attributes"`
}

type AuditLogAttributes struct {
	Action    string `json:"action"`
	ActorName string `json:"actor_name"`
	CreatedAt string `json:"created_at"`
}

type Inbox struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type EligibilitySetting struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Invitation struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes InvitationAttributes `json:"attributes"`
}

type InvitationAttributes struct {
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type CreateInvitationInput struct {
	Email    string   `json:"email"`
	GroupIDs []string `json:"group_ids,omitempty"`
}

type Group struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes GroupAttributes `json:"attributes"`
}

type GroupAttributes struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type CreateGroupInput struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions,omitempty"`
}

type UpdateGroupInput struct {
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type Member struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes MemberAttributes `json:"attributes"`
}

type MemberAttributes struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	GroupIDs []string `json:"group_ids"`
}

type UpdateMemberInput struct {
	GroupIDs []string `json:"group_ids,omitempty"`
}

func (c *Client) ListOrganizations(ctx context.Context, params PageParams) ([]Organization, error) {
	resp, err := c.Get(ctx, "/organizations", params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Organization `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetOrganizationAuditLog(ctx context.Context, orgID string, params PageParams) ([]AuditLogEntry, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/audit_log", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []AuditLogEntry `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) ListOrganizationPrograms(ctx context.Context, orgID string, params PageParams) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/programs", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) ListOrganizationInboxes(ctx context.Context, orgID string, params PageParams) ([]Inbox, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/inboxes", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Inbox `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) ListEligibilitySettings(ctx context.Context, orgID string, params PageParams) ([]EligibilitySetting, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/eligibility_settings", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []EligibilitySetting `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetEligibilitySetting(ctx context.Context, orgID, settingID string) (*EligibilitySetting, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/eligibility_settings/%s", orgID, settingID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data EligibilitySetting `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListInvitations(ctx context.Context, orgID string, params PageParams) ([]Invitation, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/invitations", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Invitation `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) CreateInvitation(ctx context.Context, orgID string, input CreateInvitationInput) (*Invitation, error) {
	body, err := wrapJSONAPI("invitation", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/organizations/%s/invitations", orgID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Invitation `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListGroups(ctx context.Context, orgID string, params PageParams) ([]Group, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/groups", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Group `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetGroup(ctx context.Context, orgID, groupID string) (*Group, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/groups/%s", orgID, groupID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Group `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateGroup(ctx context.Context, orgID string, input CreateGroupInput) (*Group, error) {
	body, err := wrapJSONAPI("group", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/organizations/%s/groups", orgID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Group `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateGroup(ctx context.Context, orgID, groupID string, input UpdateGroupInput) (*Group, error) {
	body, err := wrapJSONAPI("group", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/organizations/%s/groups/%s", orgID, groupID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Group `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListMembers(ctx context.Context, orgID string, params PageParams) ([]Member, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/members", orgID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Member `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetMember(ctx context.Context, orgID, memberID string) (*Member, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/members/%s", orgID, memberID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Member `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateMember(ctx context.Context, orgID, memberID string, input UpdateMemberInput) (*Member, error) {
	body, err := wrapJSONAPI("organization-member", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/organizations/%s/members/%s", orgID, memberID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Member `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) RemoveMember(ctx context.Context, orgID, memberID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/organizations/%s/members/%s", orgID, memberID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
