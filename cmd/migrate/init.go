package main

import (
	"context"
	"rest-fiber/config"
	"rest-fiber/internal/category"
	"rest-fiber/internal/post"
	"rest-fiber/internal/user"
 
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMigrate(ctx context.Context) error {
	logConfig := config.NewDBLogger()
	env, err := config.NewEnv()
	if err != nil {
		return err
	}

	db, err := gorm.Open(mysql.Open(env.DatabaseURL), &gorm.Config{
		Logger: logConfig,
	})
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).AutoMigrate(
		&user.User{},
		&post.Post{},
		&category.Category{},
	)
	if err != nil {
		return err
	}
	return nil
}