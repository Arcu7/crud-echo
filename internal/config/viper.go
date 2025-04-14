package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   *Server
	Database *Database
}

type Server struct {
	Host string
	Port uint16
}

type Database struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     uint16
	SSLMode  string
	TimeZone string
	LogMode  bool
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(".config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	if cfg.Database.User == "" {
		cfg.Database.User = os.Getenv("DATABASE_USERNAME")
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = os.Getenv("DATABASE_PASSWORD")
	}

	return &cfg, nil
}
