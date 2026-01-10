package post

import "context"

type PostService interface {
	FindAllPosts(ctx context.Context, page, limit, offset int) ([]PaginatedPostResponseDTO, int64, error)
	FindPostByID(ctx context.Context, id string) (*PostResponseDTO, error)
	CreatePost(ctx context.Context, dto CreatePostRequestDTO, userID string) (*PostResponseDTO, error)
	UpdatePost(ctx context.Context, id string, dto CreatePostRequestDTO, userID string) (*PostResponseDTO, error)
}
