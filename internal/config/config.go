package config

import (
	"log"
	"os"
	"time"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env string `yaml:"env"`
	HttpServer `yaml:"http_server"`
	Datasource `yaml:"datasource"`
}

type HttpServer struct {
	Address string `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	IdleTimeut time.Duration `yaml:"idle_timeout"`
}

type Datasource struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	DatabaseName string `yaml:"batabase_name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Load() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH env variable is not set")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: %s", configPath)
	}

	decoder := yaml.NewDecoder(configFile)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Error while read config file: %s", configPath)
	}

	return &cfg
}