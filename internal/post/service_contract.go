package post

import "context"

type PostService interface {
	FindAllPosts(ctx context.Context, page, limit, offset int) ([]PaginatedPostResponseDTO, int64, error)
	FindPostByID(ctx context.Context, id string) (*PostResponseDTO, error)
	CreatePost(ctx context.Context, userID string, dto CreatePostRequestDTO) (*PostResponseDTO, error)
}
