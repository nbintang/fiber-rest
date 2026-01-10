package post

import (
	"context"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]Post, int64, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, post *Post) (uuid.UUID, error)
	Update(ctx context.Context, id string, post *Post) (uuid.UUID, error)
	 Delete(ctx context.Context, id string) error
}
