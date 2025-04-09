package postgres

import (
	"crud-echo/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB() *PostgresDB {
	dsn := "user=rif password=angelcf511 dbname=project_1 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	psqldb := new(PostgresDB)
	psqldb.DB = database

	return psqldb
}

func (psqldb PostgresDB) Migrate() error {
	return psqldb.DB.AutoMigrate(&models.Books{})
}
