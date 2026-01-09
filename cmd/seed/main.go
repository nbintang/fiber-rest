package main

import (
	"flag"
	"rest-fiber/config"
	"rest-fiber/internal/infra"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	dbLogger := infra.NewDBLogger()

	env, err := config.GetEnvs()
	if err != nil {
		logrus.Fatalf("Seed failed: %v", err)
	}

	db, err := infra.GetDatabaseStandalone(env, dbLogger)
	if err != nil {
		logrus.Fatalf("Seed failed: %v", err)
	}

	countFlag := flag.String("count", "1", "specify the count")
	flag.Parse()

	count, err := strconv.Atoi(*countFlag)
	if err != nil {
		logrus.Fatalf("Invalid count: %v", err)
	}

	InitSeeds(db, Options{Count: count})
}
