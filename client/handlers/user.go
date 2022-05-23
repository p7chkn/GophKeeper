package handlers

import (
	"context"
	"fmt"
	grpc_client "new_diplom_client/grpc-client"
)

// NewUserHandler функция для создания нового обработчика пользователей
func NewUserHandler(userClient *grpc_client.UserClient) *UserHandler {
	return &UserHandler{
		UserClient: userClient,
	}
}

// UserHandler струкутра обработчика пользователей
type UserHandler struct {
	UserClient *grpc_client.UserClient
}

// RegisterUser функция регистрации пользователей
func (uh *UserHandler) RegisterUser(ctx context.Context) (string, string, error) {

	login, password, err := uh.getUserCredentials(ctx)
	if err != nil {
		return "", "", err
	}
	return uh.UserClient.Register(ctx, login, password)
}

// AuthUser функция авторизации пользователя
func (uh *UserHandler) AuthUser(ctx context.Context) (string, string, error) {
	login, password, err := uh.getUserCredentials(ctx)
	if err != nil {
		return "", "", err
	}
	return uh.UserClient.Auth(ctx, login, password)
}

func (uh *UserHandler) getUserCredentials(ctx context.Context) (string, string, error) {
	var login, password string
	fmt.Println("Enter login:")
	_, err := fmt.Scan(&login)
	if err != nil {
		fmt.Println("Error with parse login")
		return "", "", err
	}
	fmt.Println("Enter password:")
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Println("Error with parse login")
		return "", "", err
	}
	return login, password, nil
}
