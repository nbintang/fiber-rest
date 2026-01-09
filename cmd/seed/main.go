package main

import (
	"flag"
	"fmt"
	"rest-fiber/config"
	"rest-fiber/internal/infra"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
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
	countStr := *countFlag
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	InitSeeds(db, Options{
		Count: count,
	})
}
