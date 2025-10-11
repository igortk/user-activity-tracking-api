package main

import (
	log "github.com/sirupsen/logrus"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/service"
	"user-activity-tracking-api/service/database"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		log.Fatalf("Failed load configuration: %v", err)
	}

	database.ConnectDB(&cfg.DataBaseConfig)
	service.Run(&cfg.HttpConfig)
}
