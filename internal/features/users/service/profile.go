package users_service

import (
	"context"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
)

func (s *UsersService) CreateProfile(ctx context.Context, ssoUserID int64, user domain.User) (domain.User, error) {
	if ssoUserID <= 0 {
		return domain.User{}, core_errors.ErrUnauthenticated
	}
	if err := user.Validate(); err != nil {
		return domain.User{}, err
	}
	return s.usersRepository.CreateProfile(ctx, ssoUserID, user)
}

func (s *UsersService) GetProfile(ctx context.Context, ssoUserID int64) (domain.User, error) {
	if ssoUserID <= 0 {
		return domain.User{}, core_errors.ErrUnauthenticated
	}
	return s.usersRepository.GetProfile(ctx, ssoUserID)
}
