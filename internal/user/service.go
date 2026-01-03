package user

import (
	"context"
	"errors"
	"rest-fiber/pkg"

	"gorm.io/gorm"
)

type UserService interface {
	FindAllUsers(ctx context.Context) ([]UserResponse, error)
	FindUserByID(ctx context.Context, id string) (*UserResponse, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) FindAllUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return userResponses, nil
}

func (s *userService) FindUserByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrNotFound
		}
		return nil, err
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
