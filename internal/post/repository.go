package post

import (
	"context"
	"errors"
	"rest-fiber/internal/category"
	"rest-fiber/internal/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepositoryImpl{db}
}

func (r *postRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	var post Post
	var db = r.db.WithContext(ctx).Model(&post)
	err := db.Scopes(WhereID(id)).Count(&count).Error
	return count > 0, err
}

func (r *postRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]Post, int64, error) {
	var posts []Post
	var total int64
	var post Post
	var db = r.db.WithContext(ctx).Model(&post)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(
		Paginate(limit, offset),
		WhereStatus(enums.Published),
		SelectedFields(),
	).Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (r *postRepositoryImpl) FindByID(ctx context.Context, id string) (*Post, error) {
	var post Post
	if err := r.db.WithContext(ctx).
		Scopes(WhereID(id), SelectedFields("category_id")).
		Preload("Category", category.SelectedFields).
		Take(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepositoryImpl) Create(ctx context.Context, post *Post) (uuid.UUID, error) {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return uuid.Nil, err
	}
	return post.ID, nil
}

func (r *postRepositoryImpl) Update(ctx context.Context, id string, post *Post) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	tx := r.db.WithContext(ctx).
		Model(&Post{}).
		Scopes(WhereID(id)).
		Updates(post)

	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return uuid.Nil, gorm.ErrRecordNotFound
	}

	return uid, nil
}
func (r *postRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).
		Scopes(WhereID(id)).
		Delete(&Post{}).Error; err != nil {
		return err
	}
	return nil
}
