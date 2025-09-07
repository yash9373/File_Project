package database

import (
	"fmt"
	"log"
	"time"

	"file_project/config"
	"file_project/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	cfg := config.C

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.SSLMode,
	)

	gLogger := logger.New(log.New(log.Writer(), "GORM ", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
		Colorful:      true,
	})

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gLogger})
	if err != nil {
		return err
	}

	// Auto-migrate models
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}
