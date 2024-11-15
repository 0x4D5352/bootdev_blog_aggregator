package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	gatorPath, err := getConfigFilePath()
	gatorHome, err := os.ReadFile(gatorPath)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal(gatorHome, &config)
	return config
}

func getConfigFilePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	gatorPath := fmt.Sprintf("%s/%s", userHome, configFileName)
	return gatorPath, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func write(cfg *Config) error {
	contents, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	gatorPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(gatorPath, contents, 0666)
	if err != nil {
		return err
	}
	return nil
}
