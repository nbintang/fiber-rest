package post

type CreatePostRequestDTO struct {
	ImageURL   string `json:"image_url" validate:"url,required"`
	Title      string `json:"title" validate:"required,min=5,max=100"`
	Body       string `json:"body" validate:"required,min=30,max=2000"`
	Status     Status `json:"status" validate:"required,oneof=DRAFT PUBLISHED"`
	CategoryID string `json:"category_id" validate:"required,uuid4"`
}
