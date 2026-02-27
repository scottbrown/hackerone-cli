package hackeronecli

import (
	"net/url"
	"testing"
)

func TestPageParamsApply(t *testing.T) {
	p := PageParams{Number: 2, Size: 25}
	params := p.Apply(nil)

	if params.Get("page[number]") != "2" {
		t.Errorf("expected page[number]=2, got %q", params.Get("page[number]"))
	}
	if params.Get("page[size]") != "25" {
		t.Errorf("expected page[size]=25, got %q", params.Get("page[size]"))
	}
}

func TestPageParamsApplyExistingValues(t *testing.T) {
	existing := url.Values{"filter": {"active"}}
	p := PageParams{Number: 1, Size: 10}
	params := p.Apply(existing)

	if params.Get("filter") != "active" {
		t.Errorf("expected existing param preserved, got %q", params.Get("filter"))
	}
	if params.Get("page[number]") != "1" {
		t.Errorf("expected page[number]=1, got %q", params.Get("page[number]"))
	}
}

func TestPageParamsApplyZeroValues(t *testing.T) {
	p := PageParams{}
	params := p.Apply(nil)

	if params.Get("page[number]") != "" {
		t.Errorf("expected no page[number], got %q", params.Get("page[number]"))
	}
	if params.Get("page[size]") != "" {
		t.Errorf("expected no page[size], got %q", params.Get("page[size]"))
	}
}

func TestPageParamsApplyPartialValues(t *testing.T) {
	p := PageParams{Size: 50}
	params := p.Apply(nil)

	if params.Get("page[number]") != "" {
		t.Errorf("expected no page[number], got %q", params.Get("page[number]"))
	}
	if params.Get("page[size]") != "50" {
		t.Errorf("expected page[size]=50, got %q", params.Get("page[size]"))
	}
}
