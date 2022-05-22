package services

import (
	"context"
	"new_diplom/models"
)

type SecretRepoInterface interface {
	AddSecret(ctx context.Context, secret models.RawSecretData) error
	GetSecrets(ctx context.Context, userID string) ([]models.RawSecretData, error)
	DeleteSecret(ctx context.Context, secretID string, userID string) error
}

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user models.User) error
	CheckUserPassword(ctx context.Context, user models.User) (string, error)
	DeleteUser(ctx context.Context, userID string) error
}
