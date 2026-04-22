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

type Report struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes ReportAttributes `json:"attributes"`
}

type ReportAttributes struct {
	Title             string `json:"title"`
	State             string `json:"state"`
	CreatedAt         string `json:"created_at"`
	VulnerabilityInfo string `json:"vulnerability_information"`
	Severity          string `json:"severity_rating"`
	Weakness          string `json:"weakness"`
}

type CreateReportInput struct {
	Title                    string `json:"title"`
	VulnerabilityInformation string `json:"vulnerability_information"`
	Impact                   string `json:"impact,omitempty"`
	Severity                 string `json:"severity,omitempty"`
	ProgramID                string `json:"program_id"`
	WeaknessID               string `json:"weakness_id,omitempty"`
}

type CommentInput struct {
	Message  string `json:"message"`
	Internal bool   `json:"internal"`
}

type StateChangeInput struct {
	State   string `json:"state"`
	Message string `json:"message,omitempty"`
}

type SeverityInput struct {
	Rating string  `json:"rating"`
	Score  float64 `json:"score,omitempty"`
}

type BountyInput struct {
	Amount      float64 `json:"amount"`
	BonusAmount float64 `json:"bonus_amount,omitempty"`
	Message     string  `json:"message,omitempty"`
}

type BountySuggestion struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CreateBountySuggestionInput struct {
	Amount  float64 `json:"amount"`
	Message string  `json:"message,omitempty"`
}

type ListReportsFilter struct {
	Programs []string
	InboxIDs []string
}

func (c *Client) ListReports(ctx context.Context, params PageParams, filter ListReportsFilter) ([]Report, error) {
	qp := params.Apply(nil)
	for _, p := range filter.Programs {
		qp.Add("filter[program][]", p)
	}
	for _, id := range filter.InboxIDs {
		qp.Add("filter[inbox_ids][]", id)
	}
	resp, err := c.Get(ctx, "/reports", qp)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Report `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) GetReport(ctx context.Context, id string) (*Report, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/reports/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data Report `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) CreateReport(ctx context.Context, input CreateReportInput) (*Report, error) {
	body, err := wrapJSONAPI("report", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, "/reports", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data Report `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) AddComment(ctx context.Context, reportID string, input CommentInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("activity-comment", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/activities", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateAssignee(ctx context.Context, reportID string, assigneeID string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-assignee", map[string]string{"assignee_id": assigneeID})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/assignee", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CloseComments(ctx context.Context, reportID string) (map[string]interface{}, error) {
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/close_comments", reportID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCustomFields(ctx context.Context, reportID string, fields map[string]interface{}) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("custom-field-value", fields)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/custom_field_values", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCVEs(ctx context.Context, reportID string, cves []string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-cve", map[string][]string{"cves": cves})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/cves", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateInboxes(ctx context.Context, reportID string, inboxIDs []string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-inbox", map[string][]string{"inbox_ids": inboxIDs})
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/inboxes", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateSeverity(ctx context.Context, reportID string, input SeverityInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("severity", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/severities", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ChangeState(ctx context.Context, reportID string, input StateChangeInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("state-change", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/state_changes", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateReportScope(ctx context.Context, reportID string, scopeID string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("structured-scope", map[string]string{"scope_id": scopeID})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/structured_scope", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateTitle(ctx context.Context, reportID string, title string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-title", map[string]string{"title": title})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/title", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateWeakness(ctx context.Context, reportID string, weaknessID string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("weakness", map[string]string{"weakness_id": weaknessID})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/weakness", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateReference(ctx context.Context, reportID string, reference string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("issue-tracker-reference-id", map[string]string{"reference": reference})
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/issue_tracker_reference_id", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RedactReport(ctx context.Context, reportID string) (map[string]interface{}, error) {
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/redact", reportID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AddSummary(ctx context.Context, reportID string, content string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-summary", map[string]string{"content": content})
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/summaries", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GeneratePDF(ctx context.Context, reportID string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/reports/%s/pdf", reportID), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) TransferReport(ctx context.Context, reportID string, programID string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-transfer", map[string]string{"program_id": programID})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/transfer", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) EscalateReport(ctx context.Context, reportID string, integration string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-escalation", map[string]string{"integration": integration})
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/escalate", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeescalateReport(ctx context.Context, reportID string) error {
	_, err := c.Delete(ctx, fmt.Sprintf("/reports/%s/escalate", reportID))
	return err
}

func (c *Client) AddParticipant(ctx context.Context, reportID string, email string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-participant", map[string]string{"email": email})
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/participants", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UploadAttachment(ctx context.Context, reportID string, filePath string) (map[string]interface{}, error) {
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

	u := c.BaseURL + "/reports/attachments"
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

func (c *Client) DeleteAttachment(ctx context.Context, reportID string) error {
	_, err := c.Delete(ctx, fmt.Sprintf("/reports/%s/attachments", reportID))
	return err
}

func (c *Client) AwardReportBounty(ctx context.Context, reportID string, input BountyInput) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("bounty", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/bounties", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) MarkIneligible(ctx context.Context, reportID string) (map[string]interface{}, error) {
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/ineligible_for_bounty", reportID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListBountySuggestions(ctx context.Context, reportID string, params PageParams) ([]BountySuggestion, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/reports/%s/bounty_suggestions", reportID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []BountySuggestion `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func (c *Client) CreateBountySuggestion(ctx context.Context, reportID string, input CreateBountySuggestionInput) (*BountySuggestion, error) {
	body, err := wrapJSONAPI("bounty-suggestion", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/bounty_suggestions", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data BountySuggestion `json:"data"`
	}
	if err := decodeResponse(resp, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

func (c *Client) UpdateDisclosure(ctx context.Context, reportID string, state string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("disclosure-request", map[string]string{"state": state})
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/disclosure_requests", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CancelDisclosure(ctx context.Context, reportID string) error {
	_, err := c.Delete(ctx, fmt.Sprintf("/reports/%s/disclosure_requests", reportID))
	return err
}

func (c *Client) UpdateTags(ctx context.Context, reportID string, tags []string) (map[string]interface{}, error) {
	body, err := wrapJSONAPI("report-tag", map[string][]string{"tags": tags})
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/reports/%s/report_tags", reportID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RequestRetest(ctx context.Context, reportID string) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/retests", reportID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CancelRetest(ctx context.Context, reportID string) error {
	_, err := c.Delete(ctx, fmt.Sprintf("/reports/%s/retests/cancel", reportID))
	return err
}

func (c *Client) AwardSwag(ctx context.Context, reportID string) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, fmt.Sprintf("/reports/%s/swags", reportID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}
