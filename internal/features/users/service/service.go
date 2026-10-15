package users_service

import (
	"context"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateProfile(context.Context, int64, domain.User) (domain.User, error)
	GetProfile(context.Context, int64) (domain.User, error)
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
