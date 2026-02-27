package hackeronecli

import (
	"context"
	"net/url"
)

type ActivityAttributes struct {
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Activity struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes ActivityAttributes `json:"attributes"`
}

type activityResponse struct {
	Data Activity `json:"data"`
}

type activitiesResponse struct {
	Data []Activity `json:"data"`
}

func (c *Client) GetActivity(ctx context.Context, id string) (*Activity, error) {
	resp, err := c.Get(ctx, "/activities/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result activityResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListActivities(ctx context.Context, params PageParams, updatedAtAfter, updatedAtBefore string) ([]Activity, error) {
	q := params.Apply(nil)
	if q == nil {
		q = url.Values{}
	}
	if updatedAtAfter != "" {
		q.Set("page[updated_at_after]", updatedAtAfter)
	}
	if updatedAtBefore != "" {
		q.Set("page[updated_at_before]", updatedAtBefore)
	}
	resp, err := c.Get(ctx, "/incremental/activities", q)
	if err != nil {
		return nil, err
	}
	var result activitiesResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
