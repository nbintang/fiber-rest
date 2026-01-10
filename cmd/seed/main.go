package main

import (
	"flag"
	"rest-fiber/config"
	"rest-fiber/internal/infra"
	"rest-fiber/utils"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	utils.LoadEnv()
	dbLogger := infra.NewDBLogger()

	env, err := config.GetEnvs()
	if err != nil {
		logrus.Warnf("Seed failed: %v", err)
	}

	db, err := infra.GetDatabaseStandalone(env, dbLogger)
	if err != nil {
		logrus.Warnf("Seed failed: %v", err)
	}

	countFlag := flag.String("count", "1", "specify the count")
	flag.Parse()

	count, err := strconv.Atoi(*countFlag)
	if err != nil {
		logrus.Warnf("Invalid count: %v", err)
	}

	if err := InitSeeds(db, Options{Count: count}); err != nil {
		logrus.Warnf("Seed failed: %v", err)
	}
}
