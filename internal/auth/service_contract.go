package auth

import "context"

type AuthService interface {
	Register(ctx context.Context, dto *RegisterRequestDTO) error
}
