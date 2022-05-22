// Package grpc_client пакет gRPC клиентов
package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"new_diplom_client/models"
	"new_diplom_client/pb"
)

// NewSecretClient функция создания gRPC клиента для работы с секретами
func NewSecretClient(address string, access string, refresh string, userClient *UserClient) *SecretClient {
	return &SecretClient{
		address:      address,
		accessToken:  access,
		refreshToken: refresh,
		userClient:   userClient,
	}
}

type gRPCClient struct {
	pb.SecretsClient
	closeFunc func() error
}

// SecretClient структура клиента для работы с секретами
type SecretClient struct {
	address      string
	accessToken  string
	refreshToken string
	userClient   *UserClient
}

// GetSecrets функция получения секретов пользователя
func (c *SecretClient) GetSecrets(ctx context.Context) ([]models.Secret, error) {
	client, err := c.getConn()
	if err != nil {
		return nil, err
	}
	message := pb.GetSecretsRequest{}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	response, err := client.GetSecrets(ctx, &message)
	if err != nil {
		return nil, err
	}
	var result []models.Secret
	if response.Status == "unauthorized" {
		err = c.tryToRefreshToken(ctx)
		if err != nil {
			return nil, err
		}
		response, err = client.GetSecrets(ctx, &message)
		if err != nil {
			return nil, err
		}
	}
	if response.Status != "ok" {
		return nil, errors.New(response.Status)
	}
	for _, secret := range response.Secrets {
		usefulData := map[string]string{}
		for _, data := range secret.Data {
			usefulData[data.Title] = data.Value
		}
		result = append(result, models.Secret{
			ID:         secret.Id,
			Type:       secret.Type,
			MetaData:   secret.MetaData,
			UsefulData: usefulData,
		})
	}
	return result, nil
}

// CreateSecret функция создания секрета
func (c *SecretClient) CreateSecret(ctx context.Context, secret models.Secret) error {
	client, err := c.getConn()
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	message := pb.CreateSecretRequest{
		Type:     secret.Type,
		MetaData: secret.MetaData,
		Data:     secret.TransferUsefulData(),
	}
	response, err := client.CreateSecret(ctx, &message)
	if err != nil {
		return err
	}
	if response.Status == "unauthorized" {
		err = c.tryToRefreshToken(ctx)
		if err != nil {
			return err
		}
		response, err = client.CreateSecret(ctx, &message)
		if err != nil {
			return err
		}
	}
	if response.Status != "ok" {
		return errors.New(response.Status)
	}
	return nil
}

func (c *SecretClient) tryToRefreshToken(ctx context.Context) error {
	access, refresh, err := c.userClient.Refresh(ctx, c.refreshToken)
	if err != nil {
		return errors.New("failed to refresh token")
	}
	c.accessToken = access
	c.refreshToken = refresh
	return nil
}

func (c *SecretClient) getConn() (*gRPCClient, error) {
	conn, err := grpc.Dial(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := pb.NewSecretsClient(conn)

	return &gRPCClient{cl, conn.Close}, nil
}
