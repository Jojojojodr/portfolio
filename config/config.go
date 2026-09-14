package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const configFilePath = "config.yaml"

var AppConfig Config

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port  string `yaml:"port"`
		JWTSecret string `yaml:"jwt_secret"`
	} `yaml:"server"`
	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Name string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		SSLMode string `yaml:"ssl_mode"`
	} `yaml:"database"`
	Gin struct {
		Mode string `yaml:"mode"`
		TrustedProxies []string `yaml:"trusted_proxies"`
	} `yaml:"gin"`
}

func LoadConfig() {
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
}