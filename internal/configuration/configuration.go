package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Config config

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error can't get the environment variables by file")
		os.Exit(1)
	}
	if err := env.Parse(&Config); err != nil {
		logrus.Fatalf("Error initializing: %s", err.Error())
		os.Exit(1)
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
	PORT     string `env:"DB_PORT" envDefault:"3306"`
	SSLMODE  string `env:"DB_SSL_MODE" envDefault:"disable"`
	TIMEZONE string `env:"DB_TIME_ZONE"  envDefault:"America/Bogota"`
}

type APP struct {
	PORT string `env:"APP_PORT" envDefault:"8080"`
}

func (c *config) DatabaseInit() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", c.DB.USER, c.DB.PASSWORD, c.DB.HOST, c.DB.DBNAME)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
