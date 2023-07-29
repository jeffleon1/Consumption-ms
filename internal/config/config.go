package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
}

type APP struct {
	PORT string `env:"PORT" envDefault:"8080"`
}
