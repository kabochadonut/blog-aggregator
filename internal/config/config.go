package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const cfgFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) {
	cfg.CurrentUserName = username
	err := write(cfg)
	if err != nil {
		log.Fatalf("Failed to write user to file: %v", err)
	}
}

func Read() Config {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Failed to find home directory: %v", err)
	}
	content, err := os.ReadFile(cfgFilePath)
	if err != nil {
		log.Fatalf("Failed to read cfg file: %v", err)
	}
	var cfg Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse json: %v", err)
	}
	return cfg
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, cfgFileName), nil
}

func write(cfg *Config) error {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	content, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile(cfgFilePath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
