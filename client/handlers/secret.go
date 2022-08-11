package handlers

import (
	"context"
	grpc_client "new_diplom_client/grpc-client"
	"new_diplom_client/models"
)

// NewSecretHandler функция для создания нового обработчика секретов
func NewSecretHandler(client *grpc_client.SecretClient) *SecretHandler {
	return &SecretHandler{
		secretClient: client,
	}
}

// SecretHandler струкутра обработчика секретов
type SecretHandler struct {
	secretClient *grpc_client.SecretClient
	userClient   *grpc_client.UserClient
}

// CreateSecret функция создания секрета
func (sh *SecretHandler) CreateSecret(ctx context.Context, secret models.Secret) error {
	return sh.secretClient.CreateSecret(ctx, secret)
}

// GetSecret функция получения секретов
func (sh *SecretHandler) GetSecret(ctx context.Context) ([]models.Secret, error) {
	return sh.secretClient.GetSecrets(ctx)
}
