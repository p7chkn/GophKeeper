package grpc_client

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"new_diplom_client/pb"
)

// NewUserClient функция создания gRPC клиента для работы с пользователями
func NewUserClient(address string) *UserClient {
	return &UserClient{
		address: address,
	}
}

type gRPCUser struct {
	pb.UsersClient
	closeFunc func() error
}

// UserClient структура клиента для работы с пользователями
type UserClient struct {
	address string
}

// Register функция регистрации пользователя
func (u *UserClient) Register(ctx context.Context, login string, password string) (string, string, error) {
	client, err := u.getConn()
	if err != nil {
		return "", "", err
	}
	message := pb.CreateUserRequest{
		Login:    login,
		Password: password,
	}
	response, err := client.CreateUser(ctx, &message)
	if err != nil {
		return "", "", err
	}

	if response.Status == "created" {
		return response.AccessToken, response.RefreshToken, nil
	}

	return "", "", errors.New(response.Status)
}

// Auth функция авторизации пользователя
func (u *UserClient) Auth(ctx context.Context, login string, password string) (string, string, error) {
	client, err := u.getConn()
	if err != nil {
		return "", "", err
	}
	message := pb.AuthUserRequest{
		Login:    login,
		Password: password,
	}
	response, err := client.AuthUser(ctx, &message)
	if err != nil {
		return "", "", err
	}
	if response.Status == "ok" {
		return response.AccessToken, response.RefreshToken, nil
	}
	return "", "", errors.New(response.Status)
}

// Refresh функция обновления токена пользователя
func (u *UserClient) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	client, err := u.getConn()
	if err != nil {
		return "", "", err
	}
	message := pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}
	response, err := client.RefreshToken(ctx, &message)
	if err != nil {
		return "", "", nil
	}
	if response.Status == "ok" {
		return response.AccessToken, response.RefreshToken, nil
	}
	return "", "", errors.New(response.Status)
}

func (u *UserClient) getConn() (*gRPCUser, error) {
	conn, err := grpc.Dial(u.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := pb.NewUsersClient(conn)

	return &gRPCUser{cl, conn.Close}, nil
}
