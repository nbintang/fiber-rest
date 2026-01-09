package pkg

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}
}
