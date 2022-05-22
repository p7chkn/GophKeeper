package services

import (
	"context"
	"new_diplom/authorization"
	"new_diplom/models"
)

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

type UserService struct {
	db                         UserRepoInterface
	AccessTokenLiveTimeMinutes int
	RefreshTokenLiveTimeDays   int
	AccessTokenSecret          string
	RefreshTokenSecret         string
}

func (us *UserService) CreateUser(ctx context.Context, user models.User) error {
	return us.db.CreateUser(ctx, user)
}

func (us *UserService) AuthUser(ctx context.Context, user models.User) (*authorization.TokenDetails, error) {
	userID, err := us.db.CheckUserPassword(ctx, user)
	if err != nil {
		return nil, err
	}
	return authorization.CreateToken(userID, us.AccessTokenLiveTimeMinutes, us.RefreshTokenLiveTimeDays,
		us.AccessTokenSecret, us.RefreshTokenSecret)
}

func (us *UserService) DeleteUser(ctx context.Context, userID string) error {
	return us.db.DeleteUser(ctx, userID)
}

func (us *UserService) RefreshToken(ctx context.Context, refreshToken string) (*authorization.TokenDetails, error) {
	return authorization.RefreshToken(refreshToken, us.AccessTokenLiveTimeMinutes, us.RefreshTokenLiveTimeDays,
		us.AccessTokenSecret, us.RefreshTokenSecret)
}
