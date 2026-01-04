package main

import (
	"context"
	"fmt"
	"rest-fiber/config"
	"rest-fiber/internal/category"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/post"
	"rest-fiber/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitMigrate(ctx context.Context) error {
	logConfig := infra.NewDBLogger()

	env, err := config.NewEnv()
	if err != nil {
		return err
	}

	dbConf := "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC"
	dsn := fmt.Sprintf(
		dbConf,
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		env.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logConfig,
	})

	if err != nil {
		return err
	}

	if err = db.Debug().Exec(`
		DO $$ BEGIN
			CREATE TYPE role_type AS ENUM ('ADMIN', 'MEMBER');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`).Error; err != nil {
		return err
	}

	if err = db.Debug().Exec(`
		DO $$ BEGIN
			CREATE TYPE status_type AS ENUM ('PUBLISHED', 'DRAFT');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`).Error; err != nil {
		return err
	}

	return db.WithContext(ctx).Debug().AutoMigrate(
		&user.User{},
		&post.Post{},
		&category.Category{},
	)
}
