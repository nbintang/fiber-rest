package auth

import (
	"context"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/infra"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (s *authServiceImpl) generateTokens(ctx context.Context, id string, email string, role enums.EUserRoleType) (TokensResponseDto, error) {
	accessTTL := 15 * time.Minute
	refreshTTL := 24 * time.Hour
	accessJTI := uuid.NewString()
	accessExpUnix := time.Now().Add(accessTTL).Unix()
	accessToken, err := s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID: id, Email: email, Role: enums.EUserRoleType(role), JTI: accessJTI,
		Type: enums.TokenAccess}, s.env.JWTAccessSecret, accessTTL,
	)

	if err != nil {
		return TokensResponseDto{}, err
	}
	refreshJTI := uuid.NewString()
	refreshToken, err := s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID: id, Email: email, Role: enums.EUserRoleType(role), JTI: refreshJTI,
		Type: enums.TokenRefresh}, s.env.JWTRefreshSecret, refreshTTL,
	)
	if err != nil {
		return TokensResponseDto{}, err
	}

	if err := s.redisService.Set(ctx, "refresh:"+refreshJTI, id, refreshTTL); err != nil {
		return TokensResponseDto{}, err
	}
	if err := s.redisService.Set(ctx, "rt_access:"+refreshJTI, accessJTI, refreshTTL); err != nil {
		s.redisService.Del(ctx, "refresh:"+refreshJTI)
		return TokensResponseDto{}, err
	}
	if err := s.redisService.Set(ctx, "rt_access_exp:"+refreshJTI, accessExpUnix, refreshTTL); err != nil {
		s.redisService.Del(ctx, "refresh:"+refreshJTI, "rt_access:"+refreshJTI)
		return TokensResponseDto{}, err
	}
	if err := s.redisService.SAdd(ctx, "user_tokens:"+id, refreshJTI, refreshTTL); err != nil {
		s.redisService.Del(ctx,
			"refresh:"+refreshJTI,
			"rt_access:"+refreshJTI,
			"rt_access_exp:"+refreshJTI,
		)
		return TokensResponseDto{}, err
	}
	return TokensResponseDto{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authServiceImpl) generateVerificationToken(id string) (string, error) {
	return s.tokenService.GenerateToken(&infra.GenerateTokenParams{
		ID: id,
	},
		s.env.JWTVerificationSecret,
		3*time.Minute,
	)
}

func (s *authServiceImpl) revokeAllUserTokens(ctx context.Context, userID string) error {
	userTokensKey := "user_tokens:" + userID
	rtJTIs, err := s.redisService.SMembers(ctx, userTokensKey)
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	for _, rtJTI := range rtJTIs {
		s.blacklistAccessByRefreshJTI(ctx, rtJTI)
		s.redisService.Del(ctx,
			"refresh:"+rtJTI,
			"rt_access:"+rtJTI,
			"rt_access_exp:"+rtJTI,
		)
	}
	s.redisService.Del(ctx, userTokensKey)
	s.logger.Infof("revoke all tokens for user : %s", userID)
	return nil
}

func (s *authServiceImpl) blacklistAccessByRefreshJTI(ctx context.Context, rtJTI string) error {
	accessJTI, err := s.redisService.Get(ctx, "rt_access:"+rtJTI)
	if err != nil {
		if err != redis.Nil {
			s.logger.Warnf("failed get rt_access for %s: %v", rtJTI, err)
		}
		return nil
	}
	if accessJTI == "" {
		return nil
	}

	expStr, err := s.redisService.Get(ctx, "rt_access_exp:"+rtJTI)
	if err != nil {
		if err != redis.Nil {
			s.logger.Warnf("failed get rt_access_exp for %s: %v", rtJTI, err)
		}
		return nil
	}

	expUnix, _ := strconv.ParseInt(expStr, 10, 64)
	ttl := time.Until(time.Unix(expUnix, 0))
	if ttl <= 0 {
		return nil
	}

	return s.redisService.Set(ctx, "blacklist_access:"+accessJTI, "1", ttl)
}
