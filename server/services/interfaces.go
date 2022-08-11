// Package services пакет сервисов для работы с различными сущностями
package services

import (
	"context"
	"new_diplom/models"
)

// SecretRepoInterface интерфейс для взаимодействия с базой данных секретов
type SecretRepoInterface interface {
	AddSecret(ctx context.Context, secret models.RawSecretData) error
	GetSecrets(ctx context.Context, userID string) ([]models.RawSecretData, error)
	DeleteSecret(ctx context.Context, secretID string, userID string) error
}

// UserRepoInterface интерфейс для взаимодействия с базой данных пользователей
type UserRepoInterface interface {
	CreateUser(ctx context.Context, user models.User) error
	CheckUserPassword(ctx context.Context, user models.User) (string, error)
	DeleteUser(ctx context.Context, userID string) error
}
