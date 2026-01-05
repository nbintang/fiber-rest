package infra

import (
	"rest-fiber/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GenerateTokenParams struct {
	ID    string
	Email string
}

type TokenService interface {
	GenerateToken(params *GenerateTokenParams, envStr string, duration time.Duration) (string, error)
	VerifyToken(tokenStr string, envStr string) (*jwt.MapClaims, error)
}

type tokenServiceImpl struct{}

func NewTokenService(env config.Env) TokenService {
	return &tokenServiceImpl{}
}

func (s *tokenServiceImpl) GenerateToken(params *GenerateTokenParams, envStr string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    params.ID,
		"email": params.Email,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	})
	return token.SignedString([]byte(envStr))
}

func (s *tokenServiceImpl) VerifyToken(tokenStr string, envStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(envStr), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return &claims, nil
}
