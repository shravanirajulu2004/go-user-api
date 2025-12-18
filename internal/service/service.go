// internal/service/user_service.go
package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shravanirajulu2004/go-user-api/internal/models"
	"github.com/shravanirajulu2004/go-user-api/internal/repository"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id int32) (*models.UserResponse, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]models.UserResponse, int64, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("Invalid date format", zap.Error(err))
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	user, err := s.repo.CreateUser(ctx, req.Name, dob)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("User created successfully", zap.Int32("user_id", user.ID))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int32) (*models.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		s.logger.Error("Failed to get user", zap.Error(err), zap.Int32("user_id", id))
		return nil, err
	}

	age := models.CalculateAge(user.Dob)

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
		Age:  &age,
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]models.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, err := s.repo.ListUsers(ctx, int32(pageSize), int32(offset))
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.repo.CountUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to count users", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		age := models.CalculateAge(user.Dob)
		responses = append(responses, models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  user.Dob.Format("2006-01-02"),
			Age:  &age,
		})
	}

	return responses, total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("Invalid date format", zap.Error(err))
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	user, err := s.repo.UpdateUser(ctx, id, req.Name, dob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		s.logger.Error("Failed to update user", zap.Error(err), zap.Int32("user_id", id))
		return nil, err
	}

	s.logger.Info("User updated successfully", zap.Int32("user_id", user.ID))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		s.logger.Error("Failed to delete user", zap.Error(err), zap.Int32("user_id", id))
		return err
	}

	s.logger.Info("User deleted successfully", zap.Int32("user_id", id))
	return nil
}