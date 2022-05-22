// Package handlers пакет для храрнения gRPC обработчиков
package handlers

import (
	"context"
	customerrors "new_diplom/errors"
	"new_diplom/models"
	"new_diplom/pb"
)

// NewGrpcUsers функция создания обраточка запросов для пользователей
func NewGrpcUsers(userService UserServiceInterface) *GrpcUsers {
	return &GrpcUsers{
		userService: userService,
	}
}

// GrpcUsers структура для обраточика запросов для пользователя
type GrpcUsers struct {
	pb.UnimplementedUsersServer
	userService UserServiceInterface
}

// CreateUser функция создания нового пользователя
func (gh *GrpcUsers) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.TokenResponse, error) {
	user := models.User{
		Login:    in.Login,
		Password: in.Password,
	}
	err := gh.userService.CreateUser(ctx, user)
	if err != nil {
		return &pb.TokenResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	tokens, err := gh.userService.AuthUser(ctx, user)
	if err != nil {
		return &pb.TokenResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	return &pb.TokenResponse{
		Status:       "created",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// AuthUser функция авторизации пользователя
func (gh *GrpcUsers) AuthUser(ctx context.Context, in *pb.AuthUserRequest) (*pb.TokenResponse, error) {
	user := models.User{
		Login:    in.Login,
		Password: in.Password,
	}
	tokens, err := gh.userService.AuthUser(ctx, user)
	if err != nil {
		return &pb.TokenResponse{
			Status: customerrors.ParseError(err),
		}, err
	}
	return &pb.TokenResponse{
		Status:       "ok",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// RefreshToken функция для обновления токена пользователя
func (gh *GrpcUsers) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	tokens, err := gh.userService.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return &pb.TokenResponse{
			Status: customerrors.ParseError(err),
		}, err
	}
	return &pb.TokenResponse{
		Status:       "ok",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
