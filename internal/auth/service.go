package auth

import (
	"context"
	"errors"
	"rest-fiber/config"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/user"
	"rest-fiber/pkg"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authServiceImpl struct {
	userRepo     user.UserRepository
	tokenService infra.TokenService
	emailService infra.EmailService
	redisService infra.RedisService
	env          config.Env
	logger       *infra.AppLogger
}

func NewAuthService(
	userRepo user.UserRepository,
	tokenService infra.TokenService,
	emailService infra.EmailService,
	redisService infra.RedisService,
	env config.Env,
	logger *infra.AppLogger,
) AuthService {
	return &authServiceImpl{userRepo, tokenService, emailService, redisService, env, logger}
}

func (s *authServiceImpl) Register(ctx context.Context, dto *RegisterRequestDTO) error {
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
	if err != nil {
		return err
	}
	go func() {
		if err := s.emailService.SendEmail(infra.EmailParams{
			Subject: "Verification",
			Message: s.env.TargetURL + token,
			Reciever: infra.EmailReciever{
				Email: user.Email,
			}}); err != nil {
			s.logger.Error(err)
		}
	}()
	return nil
}
func (s *authServiceImpl) Login(ctx context.Context, dto *LoginRequestDTO) (Tokens, error) {
	user, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return Tokens{}, err
	}
	if user == nil {
		return Tokens{}, errors.New("User Not Found")
	}
	if err := pkg.ComparePassword(dto.Password, user.Password); err != nil {
		return Tokens{}, errors.New("Invalid Password")
	}
	if user.IsEmailVerified == false {
		return Tokens{}, errors.New("Email Not Verified")
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email)
}
func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (Tokens, error) {
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return Tokens{}, err
	}
	userID, _ := (*claims)["id"].(string)
	oldRT, ok := (*claims)["jti"].(string)
	if !ok || oldRT == "" || userID == "" {
		return Tokens{}, errors.New("invalid token claims")
	}
	storedUserID, exists, err := s.redisService.GetAndDel(ctx, "refresh:"+oldRT)
	if err != nil {
		return Tokens{}, err
	}
	if !exists || storedUserID != userID {
		s.logger.Warnf("refresh token reuse detected for user %s with JTI %s", userID, oldRT)
		s.revokeAllUserTokens(ctx, userID)
		return Tokens{}, errors.New("token reuse detected")
	}
	s.blacklistAccessByRefreshJTI(ctx, oldRT)
	s.redisService.Del(ctx, "rt_access:"+oldRT, "rt_access_exp:"+oldRT)
	s.redisService.SRem(ctx, "user_tokens:"+userID, oldRT)
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return Tokens{}, err
	}
	if u == nil {
		return Tokens{}, errors.New("User Not Found")
	}
	if !u.IsEmailVerified {
		return Tokens{}, errors.New("Email Not Verified")
	}
	return s.generateTokens(ctx, u.ID.String(), u.Email)
}


func (s *authServiceImpl) VerifyEmailToken(ctx context.Context, verificationToken string) (Tokens, error) {
	claims, err := s.tokenService.VerifyToken(verificationToken, s.env.JWTVerificationSecret)
	if err != nil {
		return Tokens{}, err
	}
	userID := (*claims)["id"].(string)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return Tokens{}, err
	}
	if user == nil {
		return Tokens{}, errors.New("User Not Found")
	}
	user.IsEmailVerified = true
	if err = s.userRepo.Update(ctx, user.ID.String(), user); err != nil {
		return Tokens{}, err
	}

	return s.generateTokens(ctx, user.ID.String(), user.Email)
}



func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("no refresh token")
	}
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return err
	}
	userID, _ := (*claims)["id"].(string)
	rtJTI, _ := (*claims)["jti"].(string)
	if userID == "" || rtJTI == "" {
		return errors.New("invalid token claims")
	}
	s.blacklistAccessByRefreshJTI(ctx, rtJTI)
	s.redisService.Del(ctx,
		"refresh:"+rtJTI,
		"rt_access:"+rtJTI,
		"rt_access_exp:"+rtJTI,
	)
	s.redisService.SRem(ctx, "user_tokens:"+userID, rtJTI)
	return nil
}

