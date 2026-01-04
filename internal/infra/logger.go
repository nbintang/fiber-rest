package infra

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type gormLogWriter struct {
	logger *logrus.Logger
}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.logger.Infof(format, args...)
}

type AppLogger struct {
	*logrus.Logger
}

func NewAppLogger() *AppLogger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.InfoLevel)
	return &AppLogger{Logger: logger}
}

type DBLogger struct {
	logger.Interface
}

func NewDBLogger(appLogger *AppLogger) *DBLogger {
	dbLogger := logger.New(
		&gormLogWriter{logger: appLogger.Logger},
		logger.Config{
			SlowThreshold:             1 * time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	return &DBLogger{Interface: dbLogger}
}
