package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env          string `yaml:"env" default:"prod"`
	Port         string `yaml:"port" default:"8080"`
	DBConnection string `yaml:"db_conn"`
}

func MustLoad(pathToConfig string) *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("Error loading .env file")
	}

	if pathToConfig == "" {
		log.Fatalf("Config file not specified")
	}

	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		log.Printf("Config file not found at path: %s", pathToConfig)
		log.Fatal("Config file does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(pathToConfig, &cfg); err != nil {
		log.Fatalf("Error reading config: %s", err.Error())
	}
	cfg.DBConnection = getDBConnection("DB_USER", "DB_PASSWORD", cfg.DBConnection)
	cfg.Port = os.Getenv("PORT")

	return &cfg
}

func getDBConnection(userEnv string, passwdEnv string, defaultConn string) string {
	user := os.Getenv(userEnv)
	passwd := os.Getenv(passwdEnv)

	if user == "" || passwd == "" {
		if defaultConn != "" {
			return defaultConn
		}
		log.Fatal("Database credentials are not set and no default connection string provided.")
	}

	res := fmt.Sprintf("mongodb://%s:%s@mongo:27017/contact?authSource=admin", user, passwd)

	return res
}
