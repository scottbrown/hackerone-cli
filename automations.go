package hackeronecli

import (
	"context"
	"fmt"
)

type Automation struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes AutomationAttributes `json:"attributes"`
}

type AutomationAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
}

type AutomationRun struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes AutomationRunAttributes `json:"attributes"`
}

type AutomationRunAttributes struct {
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type AutomationRunLog struct {
	ID         string                     `json:"id"`
	Type       string                     `json:"type"`
	Attributes AutomationRunLogAttributes `json:"attributes"`
}

type AutomationRunLogAttributes struct {
	Message   string `json:"message"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
}

func (c *Client) ListAutomations(ctx context.Context, params PageParams) ([]Automation, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, "/automations", qp)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []Automation `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetAutomation(ctx context.Context, id string) (*Automation, error) {
	resp, err := c.Get(ctx, "/automations/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data Automation `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) TriggerAutomation(ctx context.Context, id string) error {
	resp, err := c.Post(ctx, fmt.Sprintf("/automations/%s/trigger", id), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ListAutomationRuns(ctx context.Context, automationID string, params PageParams) ([]AutomationRun, error) {
	qp := params.Apply(nil)
	resp, err := c.Get(ctx, fmt.Sprintf("/automations/%s/runs", automationID), qp)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []AutomationRun `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetAutomationRun(ctx context.Context, automationID, runID string) (*AutomationRun, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/automations/%s/runs/%s", automationID, runID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data AutomationRun `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) GetAutomationRunLogs(ctx context.Context, automationID, runID string) ([]AutomationRunLog, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/automations/%s/runs/%s/logs", automationID, runID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []AutomationRunLog `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
