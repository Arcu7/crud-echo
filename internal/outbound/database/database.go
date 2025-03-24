package database

import (
	"crud-echo/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	dsn := "user=rif password=angelcf511 dbname=project_1 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return database, err
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Books{})
}
