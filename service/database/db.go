package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"user-activity-tracking-api/config"
)

var Session *gorm.DB

func ConnectDB(dbCfg *config.DataBaseConfig) {
	var err error

	Session, err = gorm.Open(postgres.Open(dbCfg.Host), &gorm.Config{})
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
	}

	log.Info("Database connection established!")
}
