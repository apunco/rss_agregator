package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"DbUrl"`
	CurrentUserName string `json:"CurrentUserName"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(c)
}

func Read() (Config, error) {
	var config Config

	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	dat, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return Config{}, err
	}

	defer dat.Close()

	fileInfo, err := os.Stat(configFilePath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		if err := createDefaultConfig(); err != nil {
			return Config{}, err
		}
		dat.Seek(0, 0)
	}

	decoder := json.NewDecoder(dat)
	if err = decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func GetConfigFilePath() (string, error) {
	if path := os.Getenv("RSS_GATOR_CONFIG"); path != "" {
		return path, nil
	}

	return filepath.Join(os.Getenv("HOME"), configFileName), nil
}

func write(cfg *Config) error {
	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	json, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, json, os.FileMode(0600))
	if err != nil {
		return err
	}

	return nil
}

func createDefaultConfig() error {
	defaultConfig := Config{
		DbUrl:           "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
		CurrentUserName: os.Getenv("USER"),
	}
	return write(&defaultConfig)

}
