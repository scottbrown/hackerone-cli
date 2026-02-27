package hackeronecli

import "context"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Website  string `json:"website"`
}

type userResponse struct {
	Data User `json:"data"`
}

func (c *Client) GetUser(ctx context.Context, username string) (*User, error) {
	resp, err := c.Get(ctx, "/users/"+username, nil)
	if err != nil {
		return nil, err
	}
	var result userResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) GetUserByID(ctx context.Context, id string) (*User, error) {
	resp, err := c.Get(ctx, "/users/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result userResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
