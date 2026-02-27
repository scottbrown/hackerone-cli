package hackeronecli

import (
	"context"
	"encoding/json"
	"net/url"
)

type AnalyticsResponse struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) GetAnalytics(ctx context.Context, program, groups, startDate, endDate string) (*AnalyticsResponse, error) {
	q := url.Values{}
	q.Set("program", program)
	if groups != "" {
		q.Set("groups", groups)
	}
	if startDate != "" {
		q.Set("start_date", startDate)
	}
	if endDate != "" {
		q.Set("end_date", endDate)
	}
	resp, err := c.Get(ctx, "/analytics", q)
	if err != nil {
		return nil, err
	}
	var result AnalyticsResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
