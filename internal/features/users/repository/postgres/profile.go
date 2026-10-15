package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/Muxammednone/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) CreateProfile(ctx context.Context, ssoUserID int64, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	// Let the unique constraint arbitrate concurrent profile creation.
	row := r.pool.QueryRow(ctx, `
 INSERT INTO todoapp.users (sso_user_id, full_name, phone_number)
 VALUES ($1, $2, $3)
 ON CONFLICT (sso_user_id) DO NOTHING
 RETURNING id, version, full_name, phone_number`, ssoUserID, user.FullName, user.PhoneNumber)
	var result domain.User
	err := row.Scan(&result.ID, &result.Version, &result.FullName, &result.PhoneNumber)
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return domain.User{}, core_errors.ErrConflict
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("create profile: %w", err)
	}
	return result, nil
}

func (r *UsersRepository) GetProfile(ctx context.Context, ssoUserID int64) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	row := r.pool.QueryRow(ctx, `
 SELECT id, version, full_name, phone_number
 FROM todoapp.users WHERE sso_user_id=$1`, ssoUserID)
	var user domain.User
	err := row.Scan(&user.ID, &user.Version, &user.FullName, &user.PhoneNumber)
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return domain.User{}, core_errors.ErrNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("get profile: %w", err)
	}
	return user, nil
}
