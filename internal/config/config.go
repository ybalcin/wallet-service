package config

import (
	"errors"
	"github.com/ybalcin/wallet-service/pkg/utility"
	"os"
)

type (
	Config struct {
		MongoSettings MongoSettings `yaml:"MongoSettings"`
	}

	MongoSettings struct {
		URI      string `yaml:"URI"`
		Database string `yaml:"Database"`
	}
)

func Read() (*Config, error) {
	mongoURI := os.Getenv("MONGO_URI")
	databaseName := os.Getenv("DATABASE_NAME")

	if utility.IsStrEmpty(mongoURI) {
		return nil, errors.New("provide valid MONGO_URI via env")
	}
	if utility.IsStrEmpty(databaseName) {
		return nil, errors.New("provide valid DATABASE_NAME via env")
	}

	return &Config{MongoSettings: MongoSettings{
		URI:      mongoURI,
		Database: databaseName,
	}}, nil
}
