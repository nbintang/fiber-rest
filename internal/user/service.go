package user

import (
	"context"
	"errors"
	"rest-fiber/internal/infra"
)

type userServiceImpl struct {
	userRepo UserRepository
	logger   *infra.AppLogger
}

func NewUserService(userRepo UserRepository, logger *infra.AppLogger) UserService {
	return &userServiceImpl{userRepo, logger}
}

func (s *userServiceImpl) FindAllUsers(ctx context.Context) ([]UserResponseDTO, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]UserResponseDTO, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, UserResponseDTO{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return userResponses, nil
}

func (s *userServiceImpl) FindUserByID(ctx context.Context, id string) (*UserResponseDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User Not Found")
	}
	return &UserResponseDTO{
		ID:        user.ID,
		AvatarURL: user.AvatarURL,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
 

