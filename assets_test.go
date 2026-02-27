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

func TestListAssets(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/assets" {
			t.Errorf("expected path /assets, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page[number]") != "1" {
			t.Errorf("expected page[number]=1, got %s", r.URL.Query().Get("page[number]"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []Asset{
				{ID: "a1", Type: "asset", Attributes: AssetAttributes{AssetType: "URL", Identifier: "example.com"}},
			},
		})
	})

	assets, err := c.ListAssets(context.Background(), PageParams{Number: 1, Size: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets))
	}
	if assets[0].ID != "a1" {
		t.Errorf("expected id a1, got %s", assets[0].ID)
	}
	if assets[0].Attributes.Identifier != "example.com" {
		t.Errorf("expected identifier example.com, got %s", assets[0].Attributes.Identifier)
	}
}

func TestGetAsset(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1" {
			t.Errorf("expected path /assets/a1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": Asset{ID: "a1", Type: "asset", Attributes: AssetAttributes{AssetType: "URL", Identifier: "example.com"}},
		})
	})

	asset, err := c.GetAsset(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.ID != "a1" {
		t.Errorf("expected id a1, got %s", asset.ID)
	}
}

func TestCreateAsset(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/assets" {
			t.Errorf("expected path /assets, got %s", r.URL.Path)
		}
		var input CreateAssetInput
		json.NewDecoder(r.Body).Decode(&input)
		if input.AssetType != "URL" {
			t.Errorf("expected asset_type URL, got %s", input.AssetType)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": Asset{ID: "a2", Type: "asset", Attributes: AssetAttributes{AssetType: "URL", Identifier: "test.com"}},
		})
	})

	asset, err := c.CreateAsset(context.Background(), CreateAssetInput{AssetType: "URL", Identifier: "test.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.ID != "a2" {
		t.Errorf("expected id a2, got %s", asset.ID)
	}
}

func TestUpdateAsset(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1" {
			t.Errorf("expected path /assets/a1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": Asset{ID: "a1", Type: "asset", Attributes: AssetAttributes{Description: "updated"}},
		})
	})

	asset, err := c.UpdateAsset(context.Background(), "a1", UpdateAssetInput{Description: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.Attributes.Description != "updated" {
		t.Errorf("expected description 'updated', got %s", asset.Attributes.Description)
	}
}

func TestArchiveAssets(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/assets" {
			t.Errorf("expected path /assets, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string][]string
		json.Unmarshal(body, &payload)
		if len(payload["ids"]) != 2 {
			t.Errorf("expected 2 ids, got %d", len(payload["ids"]))
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.ArchiveAssets(context.Background(), []string{"a1", "a2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestImportAssets(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "assets.csv")
	os.WriteFile(tmpFile, []byte("type,identifier\nURL,example.com"), 0644)

	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/assets/import" {
			t.Errorf("expected path /assets/import, got %s", r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart/form-data content-type, got %s", ct)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"import_id": "imp1"})
	})

	result, err := c.ImportAssets(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["import_id"] != "imp1" {
		t.Errorf("expected import_id imp1, got %v", result["import_id"])
	}
}

func TestGetImportStatus(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/assets/import/imp1" {
			t.Errorf("expected path /assets/import/imp1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "complete"})
	})

	result, err := c.GetImportStatus(context.Background(), "imp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "complete" {
		t.Errorf("expected status complete, got %v", result["status"])
	}
}

func TestUploadAssetScreenshot(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "screenshot.png")
	os.WriteFile(tmpFile, []byte("fake-png-data"), 0644)

	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/attachments" {
			t.Errorf("expected path /assets/a1/attachments, got %s", r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart/form-data content-type, got %s", ct)
		}
		w.WriteHeader(http.StatusCreated)
	})

	err := c.UploadAssetScreenshot(context.Background(), "a1", tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListAssetPorts(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/ports" {
			t.Errorf("expected path /assets/a1/ports, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []Port{
				{ID: "p1", Type: "port", Attributes: PortAttributes{Port: 443, Protocol: "tcp"}},
			},
		})
	})

	ports, err := c.ListAssetPorts(context.Background(), "a1", PageParams{Number: 1, Size: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(ports))
	}
	if ports[0].Attributes.Port != 443 {
		t.Errorf("expected port 443, got %d", ports[0].Attributes.Port)
	}
}

func TestCreateAssetPort(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/ports" {
			t.Errorf("expected path /assets/a1/ports, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": Port{ID: "p2", Type: "port", Attributes: PortAttributes{Port: 80, Protocol: "tcp"}},
		})
	})

	port, err := c.CreateAssetPort(context.Background(), "a1", CreatePortInput{Port: 80, Protocol: "tcp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port.ID != "p2" {
		t.Errorf("expected id p2, got %s", port.ID)
	}
}

func TestDeleteAssetPort(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/ports/p1" {
			t.Errorf("expected path /assets/a1/ports/p1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteAssetPort(context.Background(), "a1", "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetReachabilityStatus(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/reachability_status" {
			t.Errorf("expected path /assets/a1/reachability_status, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"reachable": true})
	})

	result, err := c.GetReachabilityStatus(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["reachable"] != true {
		t.Errorf("expected reachable=true, got %v", result["reachable"])
	}
}

func TestCheckReachability(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/check_reachability" {
			t.Errorf("expected path /assets/a1/check_reachability, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "checking"})
	})

	result, err := c.CheckReachability(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "checking" {
		t.Errorf("expected status checking, got %v", result["status"])
	}
}

func TestGetScannerConfig(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/scanner_configuration" {
			t.Errorf("expected path /assets/a1/scanner_configuration, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": ScannerConfiguration{Enabled: true},
		})
	})

	cfg, err := c.GetScannerConfig(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected enabled=true")
	}
}

func TestUpdateScannerConfig(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/scanner_configuration" {
			t.Errorf("expected path /assets/a1/scanner_configuration, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": ScannerConfiguration{Enabled: false},
		})
	})

	cfg, err := c.UpdateScannerConfig(context.Background(), "a1", ScannerConfiguration{Enabled: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected enabled=false")
	}
}

func TestAddAssetScope(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/assets/scopes" {
			t.Errorf("expected path /assets/scopes, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	})

	err := c.AddAssetScope(context.Background(), AssetScope{AssetID: "a1", ProgramID: "p1", Eligible: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateAssetScope(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/scopes" {
			t.Errorf("expected path /assets/a1/scopes, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.UpdateAssetScope(context.Background(), "a1", AssetScope{Eligible: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestArchiveAssetScopes(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/assets/a1/scopes" {
			t.Errorf("expected path /assets/a1/scopes, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.ArchiveAssetScopes(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListAssetTags(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/asset_tags" {
			t.Errorf("expected path /asset_tags, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []AssetTag{
				{ID: "t1", Type: "asset_tag", Attributes: AssetTagAttributes{Name: "critical"}},
			},
		})
	})

	tags, err := c.ListAssetTags(context.Background(), PageParams{Number: 1, Size: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Attributes.Name != "critical" {
		t.Errorf("expected name critical, got %s", tags[0].Attributes.Name)
	}
}

func TestCreateAssetTag(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/asset_tags" {
			t.Errorf("expected path /asset_tags, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": AssetTag{ID: "t2", Type: "asset_tag", Attributes: AssetTagAttributes{Name: "new-tag"}},
		})
	})

	tag, err := c.CreateAssetTag(context.Background(), AssetTag{Attributes: AssetTagAttributes{Name: "new-tag"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.ID != "t2" {
		t.Errorf("expected id t2, got %s", tag.ID)
	}
}

func TestListAssetTagCategories(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/asset_tag_categories" {
			t.Errorf("expected path /asset_tag_categories, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []AssetTagCategory{
				{ID: "c1", Type: "asset_tag_category", Attributes: AssetTagCategoryAttributes{Name: "severity"}},
			},
		})
	})

	cats, err := c.ListAssetTagCategories(context.Background(), PageParams{Number: 1, Size: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 1 {
		t.Fatalf("expected 1 category, got %d", len(cats))
	}
}

func TestCreateAssetTagCategory(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/asset_tag_categories" {
			t.Errorf("expected path /asset_tag_categories, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": AssetTagCategory{ID: "c2", Type: "asset_tag_category", Attributes: AssetTagCategoryAttributes{Name: "new-cat"}},
		})
	})

	cat, err := c.CreateAssetTagCategory(context.Background(), AssetTagCategory{Attributes: AssetTagCategoryAttributes{Name: "new-cat"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.ID != "c2" {
		t.Errorf("expected id c2, got %s", cat.ID)
	}
}

func TestListAssetsAPIError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
	})

	_, err := c.ListAssets(context.Background(), PageParams{})
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

func TestImportAssetsFileNotFound(t *testing.T) {
	c := NewClient("id", "tok")
	_, err := c.ImportAssets(context.Background(), "/nonexistent/file.csv")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestUploadAssetScreenshotFileNotFound(t *testing.T) {
	c := NewClient("id", "tok")
	err := c.UploadAssetScreenshot(context.Background(), "a1", "/nonexistent/file.png")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
