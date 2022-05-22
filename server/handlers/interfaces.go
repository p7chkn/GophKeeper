package handlers

import (
	"context"
	"new_diplom/authorization"
	"new_diplom/models"
)

// SecretServiceInterface интерфейс для взаимодейтсвия с сервисами секретов
type SecretServiceInterface interface {
	AddSecret(ctx context.Context, secret models.Secret) error
	GetSecrets(ctx context.Context, userID string) ([]models.SecretData, error)
	DeleteSecret(ctx context.Context, secretID string, userID string) error
}

// UserServiceInterface интерфейс для взаимодейтсвия с сервисами пользователей
type UserServiceInterface interface {
	CreateUser(ctx context.Context, user models.User) error
	AuthUser(ctx context.Context, user models.User) (*authorization.TokenDetails, error)
	DeleteUser(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*authorization.TokenDetails, error)
}
