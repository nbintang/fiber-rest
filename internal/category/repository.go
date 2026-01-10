package category

import (
	"context"

	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}


func NewCategoryRepository(db *gorm.DB) CategoryRepository{
	return &categoryRepositoryImpl{db}
}

func (r *categoryRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Category{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}