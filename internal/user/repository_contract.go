package user

import "context"

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, dto *User) error
}
