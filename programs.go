package hackeronecli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Program struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes ProgramAttributes `json:"attributes"`
}

type ProgramAttributes struct {
	Handle    string `json:"handle"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	State     string `json:"state"`
}

type ProgramBalance struct {
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

type PaymentTransaction struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CommonResponse struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Reporter struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type TeamMember struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Thanks struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Integration struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type TriageReview struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type ProgramWeakness struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CVERequest struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CreateCVERequestInput struct {
	CveIdentifier string `json:"cve_identifier,omitempty"`
	ReportID      string `json:"report_id"`
}

type HackerInvitation struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CreateHackerInvitationInput struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message,omitempty"`
}

type StructuredScope struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CreateScopeInput struct {
	AssetIdentifier       string `json:"asset_identifier"`
	AssetType             string `json:"asset_type"`
	EligibleForBounty     bool   `json:"eligible_for_bounty"`
	EligibleForSubmission bool   `json:"eligible_for_submission"`
	Instruction           string `json:"instruction,omitempty"`
}

type UpdateScopeInput struct {
	EligibleForBounty     *bool  `json:"eligible_for_bounty,omitempty"`
	EligibleForSubmission *bool  `json:"eligible_for_submission,omitempty"`
	Instruction           string `json:"instruction,omitempty"`
}

type AwardedSwag struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type UpdateSwagInput struct {
	Sent bool `json:"sent"`
}

type MessageInput struct {
	RecipientHandle string `json:"recipient_handle"`
	Message         string `json:"message"`
}

type ProgramBountyInput struct {
	ReportID    string  `json:"report_id"`
	Amount      float64 `json:"amount"`
	BonusAmount float64 `json:"bonus_amount,omitempty"`
	Message     string  `json:"message,omitempty"`
}

type PolicyInput struct {
	Policy string `json:"policy"`
}

func (c *Client) ListPrograms(ctx context.Context, params PageParams) ([]Program, error) {
	resp, err := c.Get(ctx, "/me/programs", params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Program `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) GetProgram(ctx context.Context, id string) (*Program, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data Program `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) GetProgramAuditLog(ctx context.Context, programID string, params PageParams) ([]AuditLogEntry, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/audit_log", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []AuditLogEntry `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) GetProgramBalance(ctx context.Context, programID string) (*ProgramBalance, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/balance", programID), nil)
	if err != nil {
		return nil, err
	}
	var balance ProgramBalance
	if err := decodeResponse(resp, &balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (c *Client) ListPaymentTransactions(ctx context.Context, programID string, params PageParams) ([]PaymentTransaction, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/payment_transactions", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []PaymentTransaction `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListCommonResponses(ctx context.Context, programID string, params PageParams) ([]CommonResponse, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/common_responses", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []CommonResponse `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListReporters(ctx context.Context, programID string, params PageParams) ([]Reporter, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/reporters", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Reporter `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListTeamMembers(ctx context.Context, programID string, params PageParams) ([]TeamMember, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/members", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []TeamMember `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListThanks(ctx context.Context, programID string, params PageParams) ([]Thanks, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/thanks", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Thanks `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListIntegrations(ctx context.Context, programID string, params PageParams) ([]Integration, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/integrations", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Integration `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListTriageReviews(ctx context.Context, programID string, params PageParams) ([]TriageReview, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/triage_reviews", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []TriageReview `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListWeaknesses(ctx context.Context, programID string, params PageParams) ([]ProgramWeakness, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/weaknesses", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []ProgramWeakness `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) NotifyExternalPlatform(ctx context.Context, programID string) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/notify_external_platform", programID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) SendProgramMessage(ctx context.Context, programID string, input MessageInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("program-message", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/messages", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AwardProgramBounty(ctx context.Context, programID string, input ProgramBountyInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("bounty", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/bounties", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListAllowedReporters(ctx context.Context, programID string, params PageParams) ([]Reporter, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/allowed_reporters", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Reporter `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) GetAllowedReportersHistory(ctx context.Context, programID string, params PageParams) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/allowed_reporters_history", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) GetAllowedReporterActivities(ctx context.Context, programID string, params PageParams) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/allowed_reporter_activities", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) GetAllowedReporterUsernameHistory(ctx context.Context, programID string, params PageParams) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/allowed_reporter_username_history", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) ListCVERequests(ctx context.Context, programID string, params PageParams) ([]CVERequest, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/cve_requests", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []CVERequest `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) CreateCVERequest(ctx context.Context, programID string, input CreateCVERequestInput) (*CVERequest, error) {
	body, err := wrapJSONAPI("cve-request", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/cve_requests", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data CVERequest `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) ListHackerInvitations(ctx context.Context, programID string, params PageParams) ([]HackerInvitation, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/hacker_invitations", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []HackerInvitation `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) CreateHackerInvitation(ctx context.Context, programID string, input CreateHackerInvitationInput) (*HackerInvitation, error) {
	body, err := wrapJSONAPI("hacker-invitation", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/hacker_invitations", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data HackerInvitation `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) UpdatePolicy(ctx context.Context, programID string, input PolicyInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("policy", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/programs/%s/policy", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AttachToPolicy(ctx context.Context, programID string, filePath string) (map[string]interface{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("copying file: %w", err)
	}
	writer.Close()

	u := c.BaseURL + fmt.Sprintf("/programs/%s/policy_attachments", programID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Identifier, c.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "h1-cli/"+Version)

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListScopes(ctx context.Context, programID string, params PageParams) ([]StructuredScope, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/structured_scopes", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []StructuredScope `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) CreateScope(ctx context.Context, programID string, input CreateScopeInput) (*StructuredScope, error) {
	body, err := wrapJSONAPI("structured-scope", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/structured_scopes", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data StructuredScope `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) UpdateProgramScope(ctx context.Context, programID, scopeID string, input UpdateScopeInput) (*StructuredScope, error) {
	body, err := wrapJSONAPI("structured-scope", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/programs/%s/structured_scopes/%s", programID, scopeID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data StructuredScope `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) ArchiveScope(ctx context.Context, programID, scopeID string) error {
	_, err := c.Put(ctx, fmt.Sprintf("/programs/%s/structured_scopes/%s/archive", programID, scopeID), bytes.NewReader([]byte("{}")))
	return err
}

func (c *Client) ListAwardedSwag(ctx context.Context, programID string, params PageParams) ([]AwardedSwag, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/awarded_swags", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []AwardedSwag `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) UpdateAwardedSwag(ctx context.Context, programID, swagID string, input UpdateSwagInput) (*AwardedSwag, error) {
	body, err := wrapJSONAPI("awarded-swag", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/programs/%s/awarded_swags/%s", programID, swagID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data AwardedSwag `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}
