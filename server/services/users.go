package services

import (
	"context"
	"new_diplom/authorization"
	"new_diplom/models"
)

// NewUserService функция создания сервиса для работы с пользователями
func NewUserService(db UserRepoInterface, accessTokenLiveTimeMinutes int,
	refreshTokenLiveTimeDays int, accessTokenSecret string,
	refreshTokenSecret string) *UserService {
	return &UserService{
		db:                         db,
		AccessTokenLiveTimeMinutes: accessTokenLiveTimeMinutes,
		RefreshTokenLiveTimeDays:   refreshTokenLiveTimeDays,
		AccessTokenSecret:          accessTokenSecret,
		RefreshTokenSecret:         refreshTokenSecret,
	}
}

// UserService структура для сервиса пользователей
type UserService struct {
	db                         UserRepoInterface
	AccessTokenLiveTimeMinutes int
	RefreshTokenLiveTimeDays   int
	AccessTokenSecret          string
	RefreshTokenSecret         string
}

// CreateUser функия создания пользователя
func (us *UserService) CreateUser(ctx context.Context, user models.User) error {
	return us.db.CreateUser(ctx, user)
}

// AuthUser функия авторизации пользователя
func (us *UserService) AuthUser(ctx context.Context, user models.User) (*authorization.TokenDetails, error) {
	userID, err := us.db.CheckUserPassword(ctx, user)
	if err != nil {
		return nil, err
	}
	return authorization.CreateToken(userID, us.AccessTokenLiveTimeMinutes, us.RefreshTokenLiveTimeDays,
		us.AccessTokenSecret, us.RefreshTokenSecret)
}

// DeleteUser функция удаления пользователя
func (us *UserService) DeleteUser(ctx context.Context, userID string) error {
	return us.db.DeleteUser(ctx, userID)
}

// RefreshToken функиця по обновлению токенов пользователя
func (us *UserService) RefreshToken(ctx context.Context, refreshToken string) (*authorization.TokenDetails, error) {
	return authorization.RefreshToken(refreshToken, us.AccessTokenLiveTimeMinutes, us.RefreshTokenLiveTimeDays,
		us.AccessTokenSecret, us.RefreshTokenSecret)
}
