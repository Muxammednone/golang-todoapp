package auth_client_grpc

import (
	"context"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
)

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := c.client.Login(ctx, req)
	if err != nil {
		return "", mapError(err)
	}

	return response.GetAccessToken(), nil

}
