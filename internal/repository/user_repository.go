// internal/repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/shravanirajulu2004/go-user-api/db/sqlc"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string, dob time.Time) (*sqlc.User, error)
	GetUserByID(ctx context.Context, id int32) (*sqlc.User, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]sqlc.User, error)
	UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (*sqlc.User, error)
	DeleteUser(ctx context.Context, id int32) error
	CountUsers(ctx context.Context) (int64, error)
}

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		queries: sqlc.New(db),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, name string, dob time.Time) (*sqlc.User, error) {
	user, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int32) (*sqlc.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ListUsers(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepository) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (*sqlc.User, error) {
	user, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}