package post

import (
	"context"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]Post, int64, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	Create(ctx context.Context, post *Post) (uuid.UUID, error)
}
