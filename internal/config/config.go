package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Config config

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error can't get tje environment variables by file")
	}
	if err := env.Parse(&Config); err != nil {
		logrus.Fatalf("Error initializing: %s", err.Error())
	}
}

type config struct {
	APP
	DB
}

type DB struct {
	HOST     string `env:"DB_HOST" envDefault:"localhost"`
	USER     string `env:"DB_USER" envDefault:"bia"`
	PASSWORD string `env:"DB_PASSWORD" envDefault:""`
	DBNAME   string `env:"DB_NAME" envDefault:"power_consumption"`
	PORT     string `env:"DB_PORT" envDefault:"5432"`
	SSLMODE  string `env:"DB_SSL_MODE" envDefault:"disable"`
	TIMEZONE string `env:"DB_TIME_ZONE"  envDefault:"America/Bogota"`
}

type APP struct {
	PORT string `env:"APP_PORT" envDefault:"8080"`
}

func (c *config) DatabaseInit() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", c.DB.HOST, c.DB.USER, c.DB.PASSWORD, c.DB.DBNAME, c.DB.PORT)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
