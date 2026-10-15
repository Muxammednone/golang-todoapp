package auth_client_grpc

import (
	"context"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
)

func (c *Client) Register(ctx context.Context, email, password string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	}

	response, err := c.client.Register(ctx, req)
	if err != nil {
		return 0, mapError(err)
	}

	return response.GetUserId(), nil
}
