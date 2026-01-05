package auth

import "context"

type AuthService interface {
	Register(ctx context.Context, dto *RegisterRequestDTO) error 
	VerifyEmailToken(ctx context.Context, token string) (string, error)
	Login(ctx context.Context, dto *LoginRequestDTO ) (string, error)
}
