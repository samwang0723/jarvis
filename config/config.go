package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Name string `yaml:"name"`
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Database     string `yaml:"database"`
		MaxLifetime  int    `yaml:"maxLifetime"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"database"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
}

var (
	instance Config
)

func Load() {
	f, err := os.Open("config.yaml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&instance)
	if err != nil {
		panic(err)
	}
}

func GetCurrentConfig() *Config {
	return &instance
}
