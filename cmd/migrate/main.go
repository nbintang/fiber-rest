package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	if err := InitMigrate(ctx); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
	logrus.Println("Migration Succeed")
}
