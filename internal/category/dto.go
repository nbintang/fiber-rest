package category

import "github.com/google/uuid"

type CategoryResponseDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
