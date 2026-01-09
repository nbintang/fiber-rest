package user

import (
	"context"
	"errors"
	"rest-fiber/internal/infra"
	"rest-fiber/pkg"
)

type userServiceImpl struct {
	userRepo UserRepository
	logger   *infra.AppLogger
}

func NewUserService(userRepo UserRepository, logger *infra.AppLogger) UserService {
	return &userServiceImpl{userRepo, logger}
}

func (s *userServiceImpl) FindAllUsers(ctx context.Context, page, limit, offset int) ([]UserResponseDTO, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	userResponses := pkg.MapSlices[User, UserResponseDTO](users, func(u User) UserResponseDTO {
		return UserResponseDTO{
			ID:        u.ID,
			Name:      u.Name,
			AvatarURL: u.AvatarURL,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
	})
	return userResponses, total, nil
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

func (s *userServiceImpl) UpdateProfile(ctx context.Context, id string, dto UserUpdateDTO) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	user.Name = dto.Name
	user.AvatarURL = dto.AvatarURL
	return s.userRepo.Update(ctx, id, user)
}
