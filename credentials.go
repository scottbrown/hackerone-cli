package hackeronecli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type Credential struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes CredentialAttributes `json:"attributes"`
}

type CredentialAttributes struct {
	AccountName string `json:"account_name"`
	CredType    string `json:"credential_type"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
}

type CreateCredentialInput struct {
	AccountName string `json:"account_name"`
	CredType    string `json:"credential_type"`
	Credentials string `json:"credentials"`
	ProgramID   string `json:"program_id"`
}

type UpdateCredentialInput struct {
	AccountName string `json:"account_name,omitempty"`
	Credentials string `json:"credentials,omitempty"`
}

type CredentialInquiry struct {
	ID         string                      `json:"id"`
	Type       string                      `json:"type"`
	Attributes CredentialInquiryAttributes `json:"attributes"`
}

type CredentialInquiryAttributes struct {
	Question  string `json:"question"`
	Required  bool   `json:"required"`
	FieldType string `json:"field_type"`
}

type CreateCredentialInquiryInput struct {
	Question  string `json:"question"`
	Required  bool   `json:"required"`
	FieldType string `json:"field_type"`
}

type CredentialInquiryResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (c *Client) ListCredentials(ctx context.Context, params PageParams) ([]Credential, error) {
	resp, err := c.Get(ctx, "/credentials", params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Credential `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) CreateCredential(ctx context.Context, input CreateCredentialInput) (*Credential, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, "/credentials", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Credential `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateCredential(ctx context.Context, id string, input UpdateCredentialInput) (*Credential, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/credentials/%s", id), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Credential `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) DeleteCredential(ctx context.Context, id string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/credentials/%s", id))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) AssignCredential(ctx context.Context, id string) (*Credential, error) {
	resp, err := c.Put(ctx, fmt.Sprintf("/credentials/%s/assign", id), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Credential `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) RevokeCredential(ctx context.Context, id string) (*Credential, error) {
	resp, err := c.Put(ctx, fmt.Sprintf("/credentials/%s/revoke", id), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Credential `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListCredentialInquiries(ctx context.Context, programID string, params PageParams) ([]CredentialInquiry, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/credential_inquiries", programID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []CredentialInquiry `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) CreateCredentialInquiry(ctx context.Context, programID string, input CreateCredentialInquiryInput) (*CredentialInquiry, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/programs/%s/credential_inquiries", programID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data CredentialInquiry `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateCredentialInquiry(ctx context.Context, programID, inquiryID string, input CreateCredentialInquiryInput) (*CredentialInquiry, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/programs/%s/credential_inquiries/%s", programID, inquiryID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data CredentialInquiry `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) DeleteCredentialInquiry(ctx context.Context, programID, inquiryID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/programs/%s/credential_inquiries/%s", programID, inquiryID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ListCredentialInquiryResponses(ctx context.Context, programID, inquiryID string, params PageParams) ([]CredentialInquiryResponse, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/programs/%s/credential_inquiries/%s/credential_inquiry_responses", programID, inquiryID), params.Apply(nil))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []CredentialInquiryResponse `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) DeleteCredentialInquiryResponse(ctx context.Context, programID, inquiryID, responseID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/programs/%s/credential_inquiries/%s/credential_inquiry_responses/%s", programID, inquiryID, responseID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
