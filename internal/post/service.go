package post

import (
	"context"
	"errors"
	"rest-fiber/internal/category"
	"rest-fiber/pkg/helper"
)

type postServiceImpl struct {
	postRepo     PostRepository
	categoryRepo category.CategoryRepository
}

func NewPostService(postRepo PostRepository, categoryRepo category.CategoryRepository) PostService {
	return &postServiceImpl{postRepo, categoryRepo}
}

func (s *postServiceImpl) FindAllPosts(ctx context.Context, page, limit, offset int) ([]PaginatedPostResponseDTO, int64, error) {
	posts, total, err := s.postRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	postsResponses := helper.MapSlices[Post, PaginatedPostResponseDTO](
		posts,
		func(p Post) PaginatedPostResponseDTO {
			return PaginatedPostResponseDTO{
				ID:        p.ID,
				ImageURL:  p.ImageURL,
				Title:     p.Title,
				Status:    p.Status,
				CreatedAt: p.CreatedAt,
			}
		})
	return postsResponses, total, nil
}

func (s *postServiceImpl) FindPostByID(ctx context.Context, id string) (*PostResponseDTO, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("Post Not Found")
	}
	return &PostResponseDTO{
		ID:        post.ID,
		ImageURL:  post.ImageURL,
		Title:     post.Title,
		Status:    post.Status,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Category: category.CategoryResponseDTO{
			ID:   post.Category.ID,
			Name: post.Category.Name,
		},
	}, nil
}

func (s *postServiceImpl) CreatePost(ctx context.Context, dto CreatePostRequestDTO, userID string) (*PostResponseDTO, error) {
	exists, err := s.categoryRepo.ExistsByID(ctx, dto.CategoryID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("Category Not Found")
	}

	post := &Post{
		ImageURL:   dto.ImageURL,
		Title:      dto.Title,
		Body:       dto.Body,
		CategoryID: dto.CategoryID,
		Status:     dto.Status,
		UserID:     userID,
	}

	postID, err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	created, err := s.postRepo.FindByID(ctx, postID.String())
	return &PostResponseDTO{
		ID:        created.ID,
		ImageURL:  created.ImageURL,
		Title:     created.Title,
		Status:    created.Status,
		Body:      created.Body,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
		Category: category.CategoryResponseDTO{
			ID:   created.Category.ID,
			Name: created.Category.Name,
		},
	}, nil
}

func (s *postServiceImpl) UpdatePost(ctx context.Context, id string, dto CreatePostRequestDTO, userID string) (*PostResponseDTO, error) {
	postExists, err := s.postRepo.ExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !postExists {
		return nil, errors.New("Post Not Found")
	}

	catExists, err := s.categoryRepo.ExistsByID(ctx, dto.CategoryID)
	if err != nil {
		return nil, err
	}

	if !catExists {
		return nil, errors.New("Category Not Found")
	}
	post := &Post{
		ImageURL:   dto.ImageURL,
		Title:      dto.Title,
		Body:       dto.Body,
		CategoryID: dto.CategoryID,
		Status:     dto.Status,
		UserID:     userID,
	}
	postID, err := s.postRepo.Update(ctx, id, post)
	if err != nil {
		return nil, err
	}
	updated, err := s.postRepo.FindByID(ctx, postID.String())
	return &PostResponseDTO{
		ID:        updated.ID,
		ImageURL:  updated.ImageURL,
		Title:     updated.Title,
		Status:    updated.Status,
		Body:      updated.Body,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
		Category: category.CategoryResponseDTO{
			ID:   updated.Category.ID,
			Name: updated.Category.Name,
		},
	}, nil
}
