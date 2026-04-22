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

func orgAssetsPath(orgID string) string {
	return fmt.Sprintf("/organizations/%s/assets", orgID)
}

func (c *Client) ListAssets(ctx context.Context, orgID string, params PageParams) ([]Asset, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, orgAssetsPath(orgID), qp)
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

func (c *Client) GetAsset(ctx context.Context, orgID, id string) (*Asset, error) {
	resp, err := c.Get(ctx, orgAssetsPath(orgID)+"/"+id, nil)
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

func (c *Client) CreateAsset(ctx context.Context, orgID string, input CreateAssetInput) (*Asset, error) {
	body, err := wrapJSONAPI("asset", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, orgAssetsPath(orgID), bytes.NewReader(body))
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

func (c *Client) UpdateAsset(ctx context.Context, orgID, id string, input UpdateAssetInput) (*Asset, error) {
	body, err := wrapJSONAPI("asset", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, orgAssetsPath(orgID)+"/"+id, bytes.NewReader(body))
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

func (c *Client) ArchiveAssets(ctx context.Context, orgID string, ids []string) error {
	body, err := wrapJSONAPI("asset", map[string][]string{"ids": ids})
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodDelete, orgAssetsPath(orgID), bytes.NewReader(body))
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

func (c *Client) ImportAssets(ctx context.Context, orgID, filePath string) (map[string]interface{}, error) {
	body, contentType, err := createMultipartFile("file", filePath)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, http.MethodPost, orgAssetsPath(orgID)+"/import", body)
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

func (c *Client) GetImportStatus(ctx context.Context, orgID, id string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, orgAssetsPath(orgID)+"/import/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UploadAssetScreenshot(ctx context.Context, orgID, assetID, filePath string) error {
	body, contentType, err := createMultipartFile("file", filePath)
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/attachments", orgAssetsPath(orgID), assetID), body)
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

func (c *Client) ListAssetPorts(ctx context.Context, orgID, assetID string, params PageParams) ([]Port, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, fmt.Sprintf("%s/%s/ports", orgAssetsPath(orgID), assetID), qp)
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

func (c *Client) CreateAssetPort(ctx context.Context, orgID, assetID string, input CreatePortInput) (*Port, error) {
	body, err := wrapJSONAPI("port", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("%s/%s/ports", orgAssetsPath(orgID), assetID), bytes.NewReader(body))
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

func (c *Client) DeleteAssetPort(ctx context.Context, orgID, assetID, portID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("%s/%s/ports/%s", orgAssetsPath(orgID), assetID, portID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) GetReachabilityStatus(ctx context.Context, orgID, assetID string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("%s/%s/reachability_status", orgAssetsPath(orgID), assetID), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CheckReachability(ctx context.Context, orgID, assetID string) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, fmt.Sprintf("%s/%s/check_reachability", orgAssetsPath(orgID), assetID), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetScannerConfig(ctx context.Context, orgID, assetID string) (*ScannerConfiguration, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("%s/%s/scanner_configuration", orgAssetsPath(orgID), assetID), nil)
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

func (c *Client) UpdateScannerConfig(ctx context.Context, orgID, assetID string, input ScannerConfiguration) (*ScannerConfiguration, error) {
	body, err := wrapJSONAPI("scanner-configuration", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("%s/%s/scanner_configuration", orgAssetsPath(orgID), assetID), bytes.NewReader(body))
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

func (c *Client) AddAssetScope(ctx context.Context, orgID string, input AssetScope) error {
	body, err := wrapJSONAPI("asset-scope", input)
	if err != nil {
		return err
	}
	resp, err := c.Post(ctx, orgAssetsPath(orgID)+"/scopes", bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) UpdateAssetScope(ctx context.Context, orgID, assetID string, input AssetScope) error {
	body, err := wrapJSONAPI("asset-scope", input)
	if err != nil {
		return err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("%s/%s/scopes", orgAssetsPath(orgID), assetID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ArchiveAssetScopes(ctx context.Context, orgID, assetID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("%s/%s/scopes", orgAssetsPath(orgID), assetID))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ListAssetTags(ctx context.Context, orgID string, params PageParams) ([]AssetTag, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/asset_tags", orgID), qp)
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

func (c *Client) CreateAssetTag(ctx context.Context, orgID string, input AssetTag) (*AssetTag, error) {
	body, err := wrapJSONAPI("asset-tag", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/organizations/%s/asset_tags", orgID), bytes.NewReader(body))
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

func (c *Client) ListAssetTagCategories(ctx context.Context, orgID string, params PageParams) ([]AssetTagCategory, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, fmt.Sprintf("/organizations/%s/asset_tag_categories", orgID), qp)
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

func (c *Client) CreateAssetTagCategory(ctx context.Context, orgID string, input AssetTagCategory) (*AssetTagCategory, error) {
	body, err := wrapJSONAPI("asset-tag-category", input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("/organizations/%s/asset_tag_categories", orgID), bytes.NewReader(body))
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
