package hackeronecli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Asset struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes AssetAttributes `json:"attributes"`
}

type AssetAttributes struct {
	AssetType   string `json:"asset_type"`
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
	Coverage    string `json:"coverage"`
	MaxSeverity string `json:"max_severity"`
}

type CreateAssetInput struct {
	AssetType   string `json:"asset_type"`
	Identifier  string `json:"identifier"`
	Description string `json:"description,omitempty"`
	Coverage    string `json:"coverage,omitempty"`
	MaxSeverity string `json:"max_severity,omitempty"`
}

type UpdateAssetInput struct {
	AssetType   string `json:"asset_type,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
	Description string `json:"description,omitempty"`
	Coverage    string `json:"coverage,omitempty"`
	MaxSeverity string `json:"max_severity,omitempty"`
}

type Port struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes PortAttributes `json:"attributes"`
}

type PortAttributes struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
}

type CreatePortInput struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service,omitempty"`
}

type ScannerConfiguration struct {
	Enabled bool `json:"enabled"`
}

type AssetScope struct {
	AssetID   string `json:"asset_id,omitempty"`
	ProgramID string `json:"program_id,omitempty"`
	Eligible  bool   `json:"eligible"`
}

type AssetTag struct {
	ID         string             `json:"id,omitempty"`
	Type       string             `json:"type,omitempty"`
	Attributes AssetTagAttributes `json:"attributes"`
}

type AssetTagAttributes struct {
	Name       string `json:"name"`
	CategoryID string `json:"category_id,omitempty"`
}

type AssetTagCategory struct {
	ID         string                     `json:"id,omitempty"`
	Type       string                     `json:"type,omitempty"`
	Attributes AssetTagCategoryAttributes `json:"attributes"`
}

type AssetTagCategoryAttributes struct {
	Name string `json:"name"`
}

func (c *Client) ListAssets(ctx context.Context, params PageParams) ([]Asset, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, "/assets", qp)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Asset `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetAsset(ctx context.Context, id string) (*Asset, error) {
	resp, err := c.Get(ctx, "/assets/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Asset `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateAsset(ctx context.Context, input CreateAssetInput) (*Asset, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, "/assets", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Asset `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateAsset(ctx context.Context, id string, input UpdateAssetInput) (*Asset, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, "/assets/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Asset `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ArchiveAssets(ctx context.Context, ids []string) error {
	body, err := json.Marshal(map[string][]string{"ids": ids})
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodDelete, "/assets", bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ImportAssets(ctx context.Context, filePath string) (map[string]interface{}, error) {
	body, contentType, err := createMultipartFile("file", filePath)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/assets/import", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
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

func (c *Client) GetImportStatus(ctx context.Context, id string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/assets/import/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UploadAssetScreenshot(ctx context.Context, assetID, filePath string) error {
	body, contentType, err := createMultipartFile("file", filePath)
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/assets/%s/attachments", assetID), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ListAssetPorts(ctx context.Context, assetID string, params PageParams) ([]Port, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, fmt.Sprintf("/assets/%s/ports", assetID), qp)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Port `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) CreateAssetPort(ctx context.Context, assetID string, input CreatePortInput) (*Port, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/assets/%s/ports", assetID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Port `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) DeleteAssetPort(ctx context.Context, assetID, portID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/assets/%s/ports/%s", assetID, portID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) GetReachabilityStatus(ctx context.Context, assetID string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/assets/%s/reachability_status", assetID), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CheckReachability(ctx context.Context, assetID string) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, fmt.Sprintf("/assets/%s/check_reachability", assetID), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetScannerConfig(ctx context.Context, assetID string) (*ScannerConfiguration, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/assets/%s/scanner_configuration", assetID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data ScannerConfiguration `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateScannerConfig(ctx context.Context, assetID string, input ScannerConfiguration) (*ScannerConfiguration, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/assets/%s/scanner_configuration", assetID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data ScannerConfiguration `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) AddAssetScope(ctx context.Context, input AssetScope) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := c.Post(ctx, "/assets/scopes", bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) UpdateAssetScope(ctx context.Context, assetID string, input AssetScope) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("/assets/%s/scopes", assetID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ArchiveAssetScopes(ctx context.Context, assetID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/assets/%s/scopes", assetID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ListAssetTags(ctx context.Context, params PageParams) ([]AssetTag, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, "/asset_tags", qp)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []AssetTag `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) CreateAssetTag(ctx context.Context, input AssetTag) (*AssetTag, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, "/asset_tags", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data AssetTag `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListAssetTagCategories(ctx context.Context, params PageParams) ([]AssetTagCategory, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, "/asset_tag_categories", qp)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []AssetTagCategory `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) CreateAssetTagCategory(ctx context.Context, input AssetTagCategory) (*AssetTagCategory, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, "/asset_tag_categories", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var result struct {
		Data AssetTagCategory `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func createMultipartFile(fieldName, filePath string) (io.Reader, string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return &buf, writer.FormDataContentType(), nil
}
