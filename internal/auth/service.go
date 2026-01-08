package auth

import (
	"context"
	"errors"
	"rest-fiber/config"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/user"
	"rest-fiber/pkg"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authService struct {
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
	return &authService{userRepo, tokenService, emailService, redisService, env, logger}
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
	if err != nil {
		return err
	}
	go func() {
		if err := s.emailService.SendEmail(infra.EmailParams{
			Subject: "Verification",
			Message: s.env.TargetURL + token,
			Reciever: infra.EmailReciever{
				Email: user.Email,
			},
		}); err != nil {
			s.logger.Error(err)
		}
	}()
	return nil
}

func (s *authService) Login(ctx context.Context, dto *LoginRequestDTO) (Tokens, error) {
	user, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return Tokens{}, err
	}
	if user == nil {
		return Tokens{}, errors.New("User Not Found")
	}
	if err := pkg.ComparePassword(dto.Password, user.Password); err != nil {
		s.logger.Errorf("bcrypt compare failed: %v (hash=%q)", err, user.Password)
		return Tokens{}, errors.New("Invalid Password")
	}
	if user.IsEmailVerified == false {
		return Tokens{}, errors.New("Email Not Verified")
	}

	return s.generateTokens(ctx, user.ID.String(), user.Email)
}

func (s *authService) VerifyEmailToken(ctx context.Context, verificationToken string) (Tokens, error) {
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

func (s *authService) generateVerificationToken(ID string) (string, error) {
	return s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID:   ID,
		Role: enums.Member,
	},
		s.env.JWTVerificationSecret,
		3*time.Minute,
	)
}

func (s *authService) generateTokens(ctx context.Context, ID string, Email string) (Tokens, error) {
	accessTTL := 15 * time.Minute
	refereshTTL := 24 * time.Hour
	accessJTI := uuid.NewString()
	accessExpUnix := time.Now().Add(accessTTL).Unix()
	accessToken, err := s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID:    ID,
		Email: Email,
		Role:  enums.Member,
		JTI:   accessJTI,
		Type:  enums.TokenAccess,
	},
		s.env.JWTAccessSecret,
		accessTTL,
	)
	if err != nil {
		return Tokens{}, err
	}
	refreshJTI := uuid.NewString()
	refreshToken, err := s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID:    ID,
		Email: Email,
		Role:  enums.Member,
		Type:  enums.TokenRefresh,
		JTI:   refreshJTI,
	},
		s.env.JWTRefreshSecret,
		refereshTTL,
	)
	if err != nil {
		return Tokens{}, err
	}

	refreshKey := "refresh:" + refreshJTI
	if err := s.redisService.Set(ctx, refreshKey, ID, refereshTTL); err != nil {
		return Tokens{}, err
	}

	rtAccessKey := "rt_access:" + refreshJTI
	s.redisService.Set(ctx, rtAccessKey, accessJTI, refereshTTL)
	rtAccessExpKey := "rt_access_exp:" + refreshJTI
	s.redisService.Set(ctx, rtAccessExpKey, accessExpUnix, refereshTTL)

	userTokenKey := "user_tokens:" + ID
	s.redisService.SAdd(ctx, userTokenKey, refreshJTI, refereshTTL)

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string, oldAccessToken string) (Tokens, error) {
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return Tokens{}, err
	}
	userID := (*claims)["id"].(string)
	oldRT, ok := (*claims)["jti"].(string)
	if !ok {
		return Tokens{}, errors.New("missing jti")
	}
	key := "refresh:" + oldRT
	storedUserID, exists, err := s.redisService.GetAndDel(ctx, key)
	if err != nil {
		return Tokens{}, err
	}
	if !exists || storedUserID != userID {
		s.logger.Warnf("refresh token reuse detected for user %s with JTI %s", userID, oldRT)
		s.revokeAllUserTokens(ctx, userID)
		return Tokens{}, errors.New("token reuse detected")
	}

	oldAccessJTI, err := s.redisService.Get(ctx, "rt_access:"+oldRT)
	if err == nil && oldAccessJTI != "" {
		expStr, expErr := s.redisService.Get(ctx, "rt_access_exp:"+oldRT)
		if expErr == nil {
			expUnix, _ := strconv.ParseInt(expStr, 10, 64)
			ttl := time.Until(time.Unix(expUnix, 0))
			if ttl > 0 {
				_ = s.redisService.Set(ctx, "blacklist_access:"+oldAccessJTI, "1", ttl)
			}
		}
	}

	s.redisService.Del(ctx, "rt_access:"+oldRT, "rt_access_exp:"+oldRT)
	s.redisService.SRem(ctx, "user_tokens:"+userID, oldRT)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return Tokens{}, err
	}
	if user == nil {
		return Tokens{}, errors.New("User Not Found")
	}
	if user.IsEmailVerified == false {
		return Tokens{}, errors.New("Email Not Verified")
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return err
	}
	userID := (*claims)["id"].(string)
	rtJTI := (*claims)["jti"].(string)
	accessJTI, err := s.redisService.Get(ctx, "rt_access:"+rtJTI)
	if err == nil && accessJTI != "" {
		expStr, expErr := s.redisService.Get(ctx, "rt_access_exp:"+rtJTI)
		if expErr == nil {
			expUnix, _ := strconv.ParseInt(expStr, 10, 64)
			ttl := time.Until(time.Unix(expUnix, 0))
			if ttl > 0 {
				_ = s.redisService.Set(ctx, "blacklist_access:"+accessJTI, "1", ttl)
			}
		}
	}
 
	_ = s.redisService.Del(ctx,
		"refresh:"+rtJTI,
		"rt_access:"+rtJTI,
		"rt_access_exp:"+rtJTI,
	)
	_ = s.redisService.SRem(ctx, "user_tokens:"+userID, rtJTI)
	return err
}

func (s *authService) revokeAllUserTokens(ctx context.Context, userID string) error {
	userTokensKey := "user_tokens:" + userID

	jtis, err := s.redisService.SMembers(ctx, userTokensKey)
	if err != nil {
		return err
	}

	for _, jti := range jtis {
		key := "refresh:" + jti
		s.redisService.Del(ctx, key)
	}

	s.redisService.Del(ctx, userTokensKey)

	s.logger.Infof("revoke all tokens for user : %s", userID)
	return nil
}

// func (s *authService) blacklistAccessToken(ctx context.Context, token string, ttl time.Duration) error {
// 	hash, err := pkg.HashToken(token)
// 	if err != nil {
// 		return err
// 	}
// 	key := "blacklist:" + hash
// 	return s.redis.Set(ctx, key, "1", ttl)
// }
