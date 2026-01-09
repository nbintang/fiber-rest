package main

import (
	"context"
	"rest-fiber/pkg"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	pkg.LoadEnv()
	if err := InitMigrate(ctx); err != nil {
		logrus.Warnf("Migration failed: %v", err)
	}
	logrus.Println("Migration Succeed")
}
