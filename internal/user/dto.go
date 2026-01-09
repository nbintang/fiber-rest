package user

import ( 
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
 
type UserUpdateDTO struct { 
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	AvatarURL string    `json:"avatar_url" validate:"required,url"`
}
 