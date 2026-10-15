package auth_client_grpc

import (
	"time"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
	"google.golang.org/grpc"
)

type Client struct {
	client  ssov1.SSOClient
	timeout time.Duration
}

func New(conn *grpc.ClientConn, timeout time.Duration) *Client {
	return &Client{
		client:  ssov1.NewSSOClient(conn),
		timeout: timeout,
	}
}
