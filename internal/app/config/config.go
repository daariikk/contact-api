package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env          string `yaml:"env" default:"prod"`
	Port         string `yaml:"port" default:":8080"`
	DBConnection string `yaml:"db_conn" default:"mongodb://localhost:80/contact"`
}

func MustLoad(pathToConfig string) *Config {
	if pathToConfig == "" {
		log.Fatalf("Config file not specified")
	}

	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		log.Printf("Config file not found at path: %s", pathToConfig)
		log.Fatal("Config file does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(pathToConfig, &cfg); err != nil {
		log.Fatalf("err reading config: %s", err.Error())
	}
	cfg.DBConnection = getDBConnection("DB_USER", "DB_PASSWORD", "DB_CONNECTION")

	return &cfg
}

func getDBConnection(user string, passwd string, defaultValue string) string {
	User := os.Getenv(user)
	Passwd := os.Getenv(passwd)

	if user == "" || Passwd == "" {
		return os.Getenv(defaultValue)
	}

	res := fmt.Sprintf("mongodb://%s:%s@mongo:27017/contact?authSource=admin", User, Passwd)

	return res
}
