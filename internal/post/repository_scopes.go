package post

import (
	"rest-fiber/utils/enums"

	"gorm.io/gorm"
)

type ScopeReturn func(db *gorm.DB) *gorm.DB

func SelectedFields(extra ...string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		fields := []string{
			"id",
			"image_url",
			"title",
			"body",
			"status",
			"created_at",
			"updated_at",
		}
		fields = append(fields, extra...)
		return db.Select(fields)
	}
}
func Paginate(limit, offset int) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

func WhereID(id string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
func WhereStatus(status enums.EPostStatusType) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}
