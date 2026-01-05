package user

import "gorm.io/gorm"

type ScopeReturn func(db *gorm.DB) *gorm.DB

func WhereEmail(email string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}

func WhereID(id string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func SelectPublicFields(db *gorm.DB) *gorm.DB {
	return db.Select("id", "name", "avatar_url", "password", "email", "is_email_verified")
}
