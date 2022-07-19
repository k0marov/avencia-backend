package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	FirebaseSecretPath string `json:"firebase_secret_path"`
	JWTSecretPath      string `json:"jwt_secret_path"`
	ATMSecretPath      string `json:"atm_secret_path"`
}

func LoadConfig() Config {
	configPath := os.Getenv("AVENCIA_CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("AVENCIA_CONFIG_PATH environment variable is not set. Please set it to the absolute path of the config file.")
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("error while reading the config file from path %s: %v", configPath, err)
	}
	var config Config
	json.NewDecoder(configFile).Decode(&config)
	if config.JWTSecretPath == "" || config.FirebaseSecretPath == "" || config.ATMSecretPath == "" {
		log.Fatalf("invalid config: %+v", config)
	}
	return config
}
