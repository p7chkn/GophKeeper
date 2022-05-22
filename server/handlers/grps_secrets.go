package handlers

import "C"
import (
	"context"
	"fmt"
	customerrors "new_diplom/errors"
	"new_diplom/models"
	"new_diplom/pb"
)

func NewGrpcSecrets(secretService SecretServiceInterface) *GrpcSecrets {
	return &GrpcSecrets{
		secretService: secretService,
	}
}

type GrpcSecrets struct {
	pb.UnimplementedSecretsServer
	secretService SecretServiceInterface
}

func (gh *GrpcSecrets) CreateSecret(ctx context.Context, in *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	userID := getUserFromContext(ctx)
	if userID == "" {
		return &pb.CreateSecretResponse{
			Status: "unauthorized",
		}, nil
	}

	useFullData := map[string]string{}

	for _, data := range in.Data {
		useFullData[data.Title] = data.Value
	}

	secretData := models.SecretData{
		MetaData:   in.MetaData,
		Type:       in.Type,
		UsefulData: useFullData,
	}

	secret := models.Secret{
		User: userID,
		Data: secretData,
	}

	err := gh.secretService.AddSecret(ctx, secret)
	if err != nil {
		return &pb.CreateSecretResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	return &pb.CreateSecretResponse{
		Status: "ok",
	}, nil
}

func (gh *GrpcSecrets) GetSecrets(ctx context.Context, in *pb.GetSecretsRequest) (*pb.GetSecretsResponse, error) {
	userID := getUserFromContext(ctx)
	if userID == "" {
		return &pb.GetSecretsResponse{
			Status: "unauthorized",
		}, nil
	}
	secrets, err := gh.secretService.GetSecrets(ctx, userID)
	if err != nil {
		return &pb.GetSecretsResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	var result []*pb.GetSecretsResponse_Secret
	status := "ok"
	errors := 0
	for _, secret := range secrets {
		var data []*pb.Data
		for k, v := range secret.UsefulData {
			temp := pb.Data{
				Title: k,
				Value: v,
			}
			data = append(data, &temp)
		}

		if err != nil {
			errors += 1
			status = fmt.Sprintf("parse error %v", errors)
			continue
		}
		result = append(result, &pb.GetSecretsResponse_Secret{
			Id:       secret.ID,
			MetaData: secret.MetaData,
			Type:     secret.Type,
			Data:     data,
		})
	}
	return &pb.GetSecretsResponse{
		Status:  status,
		Secrets: result,
	}, nil
}

func (gh *GrpcSecrets) DeleteSecret(ctx context.Context, in *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	err := gh.secretService.DeleteSecret(ctx, in.SecretId, "") // TODO: logic for parse user
	if err != nil {
		return &pb.DeleteSecretResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	return &pb.DeleteSecretResponse{
		Status: "ok",
	}, nil
}

func getUserFromContext(ctx context.Context) string {
	userID := ctx.Value("userID")
	if userID != nil {
		if str, ok := userID.(string); ok {
			return str
		} else {
			return ""
		}
	}
	return ""
}
