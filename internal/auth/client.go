package auth

import (
	"context"
	"fmt"
	"time"

	pb "github.com/DailyPepper/auth-service/pkg/generated/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	client  pb.AuthServiceClient
	timeout time.Duration
}

func NewAuthClient(authServiceURL string) (*AuthClient, error) {
	conn, err := grpc.Dial(authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		conn:    conn,
		client:  pb.NewAuthServiceClient(conn),
		timeout: 10 * time.Second,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) Register(ctx context.Context, email, password, firstName, lastName string) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.Register(ctx, &pb.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		Surname:   lastName,
	})
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*pb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
}
