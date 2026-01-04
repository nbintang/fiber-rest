package infra

import (
	"context"
	"fmt"
	"rest-fiber/config"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(lc fx.Lifecycle, env config.Env, logger *DBLogger) (*gorm.DB, error) {
	dbConf := "host=%s user=%s password=%s dbname=%s port=%d"
	dsn := fmt.Sprintf(
		dbConf,
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		env.DatabasePort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return sqlDB.Close()
		},
		OnStart: func(ctx context.Context) error {
			return nil
		},
	})

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
