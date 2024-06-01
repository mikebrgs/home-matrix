package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Dbname     string `yaml:"dbname"`
		DisableSSL bool   `yaml:"disablessl"`
	} `yaml:"database"`
	MQTT struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"mqtt"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
