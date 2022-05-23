package services

import (
	"context"
	"new_diplom/models"
)

// NewSecretService функция создания сервиса для работы с секретами
func NewSecretService(db SecretRepoInterface) *SecretServices {
	return &SecretServices{
		db: db,
	}
}

// SecretServices структура сервиса секретов
type SecretServices struct {
	db SecretRepoInterface
}

// AddSecret функция по созданию нового секрета
func (ss *SecretServices) AddSecret(ctx context.Context, secret models.Secret) error {
	rawSecret, err := models.NewRawSecretData(secret)
	if err != nil {
		return err
	}
	return ss.db.AddSecret(ctx, *rawSecret)
}

// GetSecrets функия по получению всех секретов пользователя
func (ss *SecretServices) GetSecrets(ctx context.Context, userID string) ([]models.SecretData, error) {
	rawSecrets, err := ss.db.GetSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []models.SecretData
	for _, secret := range rawSecrets {
		data, err := secret.DecryptToSecretData()
		if err != nil {
			continue
		}
		result = append(result, *data)
	}
	return result, nil
}

// DeleteSecret функция удаления скрета
func (ss *SecretServices) DeleteSecret(ctx context.Context, secretID string,
	userID string) error {
	return ss.db.DeleteSecret(ctx, secretID, userID)
}
