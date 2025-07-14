package dbx

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"heart/internal/config"
)

var DB *gorm.DB

func Connect(cfg *config.DBConfig) error {
	dsn := fmt.Sprintf("host=%s user='%s' password='%s' dbname=%s port=%s sslmode=%s application_name='%s'",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.AppName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}
