package post

import (
	"database/sql/driver"
	"fmt"
	"rest-fiber/internal/category"
	"rest-fiber/utils/enums"
	"time"

	"github.com/google/uuid"
)

type Status enums.EPostStatusType

type Post struct {
	ID         uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ImageURL   string            `gorm:"type:text;null;default:null;column:image_url"`
	Title      string            `gorm:"type:varchar(255);not null;column:title"`
	Body       string            `gorm:"type:text;not null;column:body"`
	UserID     string            `gorm:"type:char(36);not null;column:user_id"`
	CategoryID string            `gorm:"type:char(36);not null;column:category_id"`
	Status     Status            `gorm:"type:status_type;not null;default:'DRAFT';column:status"`
	Category   category.Category `gorm:"foreignKey:CategoryID;references:ID"`
	CreatedAt  time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time         `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  *time.Time       `gorm:"index"` 
}

func (p *Post) TableName() string {
	return "posts"
}

func (p *Post) IsPublished() bool {
	return p.Status == Status(enums.Published)
}

func (r *Status) Scan(value any) error {
	if value == nil {
		*r = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*r = Status(string(v))
	case string:
		*r = Status(v)
	default:
		return fmt.Errorf("cannot scan %T into Role", value)
	}
	return nil
}

func (r Status) Value() (driver.Value, error) {
	return string(r), nil
}
