package auth

import (
	"context"
	"errors"
	"fmt"
	"rest-fiber/config"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/user"
	"rest-fiber/pkg"
	"time"
)

type authService struct {
	userRepo     user.UserRepository
	tokenService infra.TokenService
	emailService infra.EmailService
	env          config.Env
	logger       *infra.AppLogger
}

func NewAuthService(
	userRepo user.UserRepository,
	tokenService infra.TokenService,
	emailService infra.EmailService,
	env config.Env,
	logger *infra.AppLogger,
) AuthService {
	return &authService{userRepo, tokenService, emailService, env, logger}
}

func (s *authService) Register(ctx context.Context, dto *RegisterRequestDTO) error {
	exists, err := s.userRepo.FindExistsByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("User Already Exist")
	}
	hashed, err := pkg.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	user := user.User{
		AvatarURL: dto.AvatarURL,
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  hashed,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return err
	}

	token, err := s.generateVerificationToken(user.ID.String())
	frontendURL := "http://localhost:8080/api/auth/verify?token=%s"
	if err != nil {
		return err
	}
	go func() {
		if err := s.emailService.SendEmail(infra.EmailParams{
			Subject: "Verification",
			Message: fmt.Sprintf(frontendURL, token),
			Reciever: infra.EmailReciever{
				Email: user.Email,
			},
		}); err != nil {
			s.logger.Error(err)
		}
	}()

	return nil
}

func (s *authService) Login(ctx context.Context, dto *LoginRequestDTO) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("User Not Found")
	}
	if err := pkg.ComparePassword(dto.Password, user.Password); err != nil {
    s.logger.Errorf("bcrypt compare failed: %v (hash=%q)", err, user.Password)
    return "", errors.New("Invalid Password")
}
	if user.IsEmailVerified == false {
		return "", errors.New("Email Not Verified")
	}
	return s.generateAccessToken(user.ID.String(), user.Email)
}

func (s *authService) VerifyEmailToken(ctx context.Context, token string) (string, error) {
	claims, err := s.tokenService.VerifyToken(token, s.env.JWTVerificationSecret)
	if err != nil {
		return "", err
	}
	userID := (*claims)["id"].(string)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("User Not Found")
	}
	user.IsEmailVerified = true
	if err = s.userRepo.Update(ctx, user.ID.String(), user); err != nil {
		return "", err
	}

	return s.generateAccessToken(user.ID.String(), user.Email)
}


func (s *authService) generateVerificationToken(ID string) (string, error) {
	return s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID: ID,
	},
		s.env.JWTVerificationSecret,
		3*time.Minute,
	)
}

func (a *authService) generateAccessToken(ID string, Email string) (string, error) {
	return a.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID: ID,
	},
		a.env.JWTAccessSecret,
		3*time.Minute,
	)
}

