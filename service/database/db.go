package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"user-activity-tracking-api/config"
)

type Client struct {
	db *gorm.DB
}

func NewClient(dbCfg *config.DataBaseConfig) *Client {
	cl := &Client{}
	cl.db = cl.connectDB(dbCfg)

	return cl
}

func (c *Client) GetDb() *gorm.DB {
	return c.db
}

func (c *Client) Close() {
	db, _ := c.db.DB()
	db.Close()
}

func (c *Client) connectDB(dbCfg *config.DataBaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbCfg.Host), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying *sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(int(dbCfg.MaxOpenConns))
	sqlDB.SetMaxIdleConns(int(dbCfg.MaxIdleConns))
	sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.ConnMaxLifetime) * time.Minute)

	log.Info("Database connection established!")
	return db
}
