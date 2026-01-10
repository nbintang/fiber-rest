package main

import (
	"context"
	"rest-fiber/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	utils.LoadEnv()
	if err := InitMigrate(ctx); err != nil {
		logrus.Warnf("Migration failed: %v", err)
	}
	logrus.Println("Migration Succeed")
}
