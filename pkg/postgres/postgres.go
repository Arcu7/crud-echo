package postgres

import (
	"crud-echo/internal/config"
	"crud-echo/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewDB(cfg *config.Config) (*PostgresDB, error) {
	dsn := "user=rif password=angelcf511 dbname=project_1 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn = fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		// NOTE: don't know why, but this will only work if the host is not specified
		// cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	database, err := gorm.Open(postgres.Open(dsn))
	if cfg.Database.LogMode {
		database.Logger = logger.Default.LogMode(logger.Info)
	} else {
		database.Logger = logger.Default.LogMode(logger.Silent)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	psqldb := new(PostgresDB)
	psqldb.DB = database

	return psqldb, nil
}

func (psqldb *PostgresDB) Migrate() error {
	return psqldb.DB.AutoMigrate(&models.Books{})
}

func (psqldb *PostgresDB) GetDB() *gorm.DB {
	return psqldb.DB
}
