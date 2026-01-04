package user

import "context"

type UserService interface {
	FindAllUsers(ctx context.Context) ([]UserResponseDTO, error)
	FindUserByID(ctx context.Context, id string) (*UserResponseDTO, error)
}
