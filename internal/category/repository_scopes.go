package category

import "gorm.io/gorm"

func SelectPublicFields(db *gorm.DB) *gorm.DB{
	return db.Select("id", "name");
}