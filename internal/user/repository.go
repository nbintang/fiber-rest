package user

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id string) (User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context) ([]User, error) {
	var user []User
	err := r.db.WithContext(ctx).Find(&user).Error
	return user, err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).Take(&user, "id = ?", id).Error
	return user, err
}
