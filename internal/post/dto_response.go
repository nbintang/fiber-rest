package post

import (
	"rest-fiber/internal/category"
	"time"

	"github.com/google/uuid"
)

type PaginatedPostResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type PostResponseDTO struct {
	ID        uuid.UUID                    `json:"id"`
	ImageURL  string                       `json:"image_url"`
	Title     string                       `json:"title"`
	Body      string                       `json:"body"`
	Status    Status                       `json:"status"`
	Category  category.CategoryResponseDTO `json:"category"`
	CreatedAt time.Time                    `json:"created_at"`
	UpdatedAt time.Time                    `json:"updated_at"`
}
