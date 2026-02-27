package hackeronecli

import (
	"bytes"
	"context"
	"encoding/json"
)

type SendEmailInput struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (c *Client) SendEmail(ctx context.Context, input SendEmailInput) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := c.Post(ctx, "/email", bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
