package config

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

type Config struct {
	FirebaseSecretPath string `json:"firebase_secret_path"`
	JWTSecretPath      string `json:"jwt_secret_path"`
	ATMSecretPath      string `json:"atm_secret_path"`

  StaticDir string `json:"static_dir"`
  StaticHost string `json:"static_host"`
}

func LoadConfig() Config {
	configPath := os.Getenv("AVENCIA_CONFIG_PATH")
	if configPath == "" {
		log.Fatalf(`AVENCIA_CONFIG_PATH environment variable is not set. 
		Please set it to the absolute path of the config file.`)
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("error while reading the config file from path %s: %v", configPath, err)
	}
	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatalf("error while parsing the config file %s: %v", configPath, err)
	}
	checkConfigFields(config)
	return config
}

func checkConfigFields(conf Config) {
  rConf := reflect.ValueOf(&conf)	.Elem()
  for i := 0; i < rConf.NumField(); i++ {
  	f := rConf.Field(i)
  	if f.IsZero() {
			log.Fatalf("invalid config: field %s is unset", rConf.Type().Field(i).Tag.Get("json"))
  	}
  }
}
