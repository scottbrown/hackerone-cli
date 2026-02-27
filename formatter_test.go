package hackeronecli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"text", FormatText, false},
		{"TEXT", FormatText, false},
		{"markdown", FormatMarkdown, false},
		{"MARKDOWN", FormatMarkdown, false},
		{"xml", "", true},
		{"", "", true},
		{"csv", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

type testAttributes struct {
	Title     string `json:"title"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
}

type testItem struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes testAttributes `json:"attributes"`
}

type testFlat struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type testMapAttributes struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

func TestFormatOutputJSONSingle(t *testing.T) {
	item := testItem{
		ID:   "123",
		Type: "report",
		Attributes: testAttributes{
			Title:     "XSS bug",
			State:     "triaged",
			CreatedAt: "2024-06-15",
		},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatJSON, item); err != nil {
		t.Fatal(err)
	}
	var decoded testItem
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded.ID != "123" {
		t.Errorf("expected id 123, got %s", decoded.ID)
	}
	if !strings.Contains(buf.String(), "  ") {
		t.Error("expected indented JSON output")
	}
}

func TestFormatOutputJSONSlice(t *testing.T) {
	items := []testFlat{
		{ID: "1", Username: "alice", Name: "Alice"},
		{ID: "2", Username: "bob", Name: "Bob"},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatJSON, items); err != nil {
		t.Fatal(err)
	}
	var decoded []testFlat
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded) != 2 {
		t.Errorf("expected 2 items, got %d", len(decoded))
	}
}

func TestFormatOutputTextSingle(t *testing.T) {
	item := testFlat{ID: "42", Username: "alice", Name: "Alice Smith"}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, item); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "id: 42") {
		t.Errorf("expected 'id: 42' in output, got: %s", out)
	}
	if !strings.Contains(out, "username: alice") {
		t.Errorf("expected 'username: alice' in output, got: %s", out)
	}
	if !strings.Contains(out, "name: Alice Smith") {
		t.Errorf("expected 'name: Alice Smith' in output, got: %s", out)
	}
}

func TestFormatOutputTextSingleWithAttributes(t *testing.T) {
	item := testItem{
		ID:   "123",
		Type: "report",
		Attributes: testAttributes{
			Title:     "XSS bug",
			State:     "triaged",
			CreatedAt: "2024-06-15",
		},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, item); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "id: 123") {
		t.Errorf("expected 'id: 123', got: %s", out)
	}
	if !strings.Contains(out, "title: XSS bug") {
		t.Errorf("expected flattened 'title: XSS bug', got: %s", out)
	}
	if !strings.Contains(out, "state: triaged") {
		t.Errorf("expected flattened 'state: triaged', got: %s", out)
	}
	if strings.Contains(out, "attributes:") {
		t.Errorf("Attributes should be flattened, not shown as a key")
	}
}

func TestFormatOutputTextSlice(t *testing.T) {
	items := []testItem{
		{
			ID:   "1",
			Type: "report",
			Attributes: testAttributes{
				Title: "XSS",
				State: "triaged",
			},
		},
		{
			ID:   "2",
			Type: "report",
			Attributes: testAttributes{
				Title: "CSRF",
				State: "new",
			},
		},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, items); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d: %s", len(lines), out)
	}
	if !strings.Contains(lines[0], "id") {
		t.Errorf("expected header to contain 'id', got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "title") {
		t.Errorf("expected header to contain 'title' (flattened), got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "XSS") {
		t.Errorf("expected first row to contain 'XSS', got: %s", lines[1])
	}
}

func TestFormatOutputTextMap(t *testing.T) {
	m := map[string]interface{}{
		"beta":  "two",
		"alpha": "one",
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, m); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "alpha:") {
		t.Errorf("expected sorted keys, first should be alpha, got: %s", lines[0])
	}
}

func TestFormatOutputTextString(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, "hello world"); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "hello world" {
		t.Errorf("expected 'hello world', got: %s", buf.String())
	}
}

func TestFormatOutputTextPointer(t *testing.T) {
	item := &testFlat{ID: "10", Username: "ptr", Name: "Pointer"}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, item); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "id: 10") {
		t.Errorf("expected 'id: 10', got: %s", buf.String())
	}
}

func TestFormatOutputTextNilPointer(t *testing.T) {
	var item *testFlat
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, item); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "<nil>") {
		t.Errorf("expected '<nil>', got: %s", buf.String())
	}
}

func TestFormatOutputTextEmptySlice(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, []testFlat{}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output for empty slice, got: %q", buf.String())
	}
}

func TestFormatOutputTextMapAttributes(t *testing.T) {
	item := testMapAttributes{
		ID:   "99",
		Type: "transaction",
		Attributes: map[string]interface{}{
			"amount":   100.5,
			"currency": "USD",
		},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, item); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "id: 99") {
		t.Errorf("expected 'id: 99', got: %s", out)
	}
	if !strings.Contains(out, "amount:") {
		t.Errorf("expected flattened 'amount' from map attributes, got: %s", out)
	}
	if !strings.Contains(out, "currency:") {
		t.Errorf("expected flattened 'currency' from map attributes, got: %s", out)
	}
}

func TestFormatOutputMarkdownSingle(t *testing.T) {
	item := testFlat{ID: "42", Username: "alice", Name: "Alice Smith"}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, item); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key | Value |") {
		t.Errorf("expected markdown header row, got: %s", out)
	}
	if !strings.Contains(out, "| --- | --- |") {
		t.Errorf("expected markdown separator, got: %s", out)
	}
	if !strings.Contains(out, "| id | 42 |") {
		t.Errorf("expected '| id | 42 |', got: %s", out)
	}
}

func TestFormatOutputMarkdownSingleWithAttributes(t *testing.T) {
	item := testItem{
		ID:   "123",
		Type: "report",
		Attributes: testAttributes{
			Title: "XSS bug",
			State: "triaged",
		},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, item); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "| title | XSS bug |") {
		t.Errorf("expected flattened attributes in markdown, got: %s", out)
	}
}

func TestFormatOutputMarkdownSlice(t *testing.T) {
	items := []testItem{
		{
			ID:   "1",
			Type: "report",
			Attributes: testAttributes{
				Title: "XSS",
				State: "triaged",
			},
		},
		{
			ID:   "2",
			Type: "report",
			Attributes: testAttributes{
				Title: "CSRF",
				State: "new",
			},
		},
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, items); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (header + sep + 2 rows), got %d: %s", len(lines), out)
	}
	if !strings.Contains(lines[0], "| id |") {
		t.Errorf("expected header with '| id |', got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "| --- |") {
		t.Errorf("expected separator, got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "XSS") {
		t.Errorf("expected row with XSS, got: %s", lines[2])
	}
}

func TestFormatOutputMarkdownMap(t *testing.T) {
	m := map[string]interface{}{
		"beta":  "two",
		"alpha": "one",
	}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, m); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key | Value |") {
		t.Errorf("expected markdown key-value table, got: %s", out)
	}
	if !strings.Contains(out, "| alpha | one |") {
		t.Errorf("expected sorted key alpha first, got: %s", out)
	}
}

func TestFormatOutputMarkdownEmptySlice(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, []testFlat{}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output for empty slice, got: %q", buf.String())
	}
}

func TestFormatOutputMarkdownPointer(t *testing.T) {
	item := &testFlat{ID: "10", Username: "ptr", Name: "Pointer"}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, item); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "| id | 10 |") {
		t.Errorf("expected '| id | 10 |', got: %s", buf.String())
	}
}

func TestFormatOutputMarkdownNilPointer(t *testing.T) {
	var item *testFlat
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, item); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "*nil*") {
		t.Errorf("expected '*nil*', got: %s", buf.String())
	}
}

func TestFormatOutputMarkdownString(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, "hello world"); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "hello world" {
		t.Errorf("expected 'hello world', got: %s", buf.String())
	}
}

func TestFormatMessageJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatMessage(&buf, FormatJSON, "done"); err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["message"] != "done" {
		t.Errorf("expected message 'done', got %q", m["message"])
	}
}

func TestFormatMessageText(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatMessage(&buf, FormatText, "done"); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "done" {
		t.Errorf("expected 'done', got %q", buf.String())
	}
}

func TestFormatMessageMarkdown(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatMessage(&buf, FormatMarkdown, "done"); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "done" {
		t.Errorf("expected 'done', got %q", buf.String())
	}
}

func TestFormatOutputTextSliceOfStrings(t *testing.T) {
	items := []interface{}{"one", "two", "three"}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatText, items); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %s", len(lines), buf.String())
	}
}

func TestFormatOutputMarkdownSliceOfStrings(t *testing.T) {
	items := []interface{}{"one", "two", "three"}
	var buf bytes.Buffer
	if err := FormatOutput(&buf, FormatMarkdown, items); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %s", len(lines), buf.String())
	}
}

func TestFormatOutputWithRealTypes(t *testing.T) {
	report := Report{
		ID:   "12345",
		Type: "report",
		Attributes: ReportAttributes{
			Title:    "XSS in login page",
			State:    "triaged",
			Severity: "high",
		},
	}

	t.Run("text single report", func(t *testing.T) {
		var buf bytes.Buffer
		if err := FormatOutput(&buf, FormatText, report); err != nil {
			t.Fatal(err)
		}
		out := buf.String()
		if !strings.Contains(out, "id: 12345") {
			t.Errorf("expected 'id: 12345', got: %s", out)
		}
		if !strings.Contains(out, "title: XSS in login page") {
			t.Errorf("expected flattened title, got: %s", out)
		}
	})

	t.Run("text report slice", func(t *testing.T) {
		reports := []Report{report, {
			ID:   "12346",
			Type: "report",
			Attributes: ReportAttributes{
				Title: "CSRF on settings",
				State: "new",
			},
		}}
		var buf bytes.Buffer
		if err := FormatOutput(&buf, FormatText, reports); err != nil {
			t.Fatal(err)
		}
		out := buf.String()
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d", len(lines))
		}
	})

	t.Run("markdown report slice", func(t *testing.T) {
		reports := []Report{report}
		var buf bytes.Buffer
		if err := FormatOutput(&buf, FormatMarkdown, reports); err != nil {
			t.Fatal(err)
		}
		out := buf.String()
		if !strings.Contains(out, "| id | type | title") {
			t.Errorf("expected markdown header with flattened columns, got: %s", out)
		}
	})
}

func TestFormatOutputMapSlice(t *testing.T) {
	items := []map[string]interface{}{
		{"id": "1", "name": "alpha"},
		{"id": "2", "name": "beta"},
	}
	t.Run("text", func(t *testing.T) {
		var buf bytes.Buffer
		if err := FormatOutput(&buf, FormatText, items); err != nil {
			t.Fatal(err)
		}
		out := buf.String()
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines (header + 2 rows), got %d: %s", len(lines), out)
		}
	})
	t.Run("markdown", func(t *testing.T) {
		var buf bytes.Buffer
		if err := FormatOutput(&buf, FormatMarkdown, items); err != nil {
			t.Fatal(err)
		}
		out := buf.String()
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 4 {
			t.Fatalf("expected 4 lines (header + sep + 2 rows), got %d: %s", len(lines), out)
		}
	})
}
