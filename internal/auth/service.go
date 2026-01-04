package auth

import (
	"context"
	"rest-fiber/internal/user"
	"rest-fiber/pkg"
)

type authService struct {
	userRepo user.UserRepository
}

func NewAuthService(userRepo user.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, dto *RegisterRequestDTO) error {
	exists, err := s.userRepo.FindExistsByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}
	if exists {
		return pkg.ErrAlreadyExists
	}
	hashed, err := pkg.HashPassword(dto.Password)
	if err != nil {
		return err
	}
	user := &user.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashed,
	}
	return s.userRepo.Create(ctx, user)
}
