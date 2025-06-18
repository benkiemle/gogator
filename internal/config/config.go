package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CurrentUserName  string `json:"current_user_name"`
	ConnectionString string `json:"connection_string"`
}

func getGatorConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.gatorconfig.json", nil
}

func Read() (*Config, error) {
	filePath, err := getGatorConfigPath()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (config *Config) SetUser(username string) error {
	filePath, err := getGatorConfigPath()
	if err != nil {
		return err
	}
	config.CurrentUserName = username

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(config)
	return nil
}
