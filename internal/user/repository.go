package user

import (
	"context"
	"errors"
	"rest-fiber/internal/enums"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]User, int64, error) {
	var users []User
	var total int64
	var user User
	db := r.db.WithContext(ctx).Model(&user)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(Paginate(limit, offset)).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Scopes(WhereID(id), SelectPublicFields).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Scopes(WhereEmail(email), SelectPublicFields).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	var user User
	err := r.db.WithContext(ctx).Model(&user).Scopes(WhereEmail(email)).Count(&count).Error
	return count > 0, err
}

func (r *userRepositoryImpl) Update(ctx context.Context, id string, dto *User) error {
	if err := r.db.WithContext(ctx).Scopes(WhereID(id)).Updates(&dto).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, dto *User) error {
	dto.Role = Role(enums.Member)
	err := r.db.WithContext(ctx).Create(&dto).Error
	return err
}
