package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sano/utils"
)

type Cfg struct {
	App         App         `yaml:"app"`
	HealthCheck HealthCheck `yaml:"healthCheck"`
	Services    Services    `yaml:"services"`
}

type Services []Service

type App struct {
	Port     int      `yaml:"port"`
	Database Database `yaml:"database"`
}

type Database struct {
	Url string `yaml:"url"`
}

type HealthCheck struct {
	Cron string `yaml:"cron"`
}

type Service struct {
	Name        string  `yaml:"name"`
	DisplayName *string `yaml:"displayName"`
	Url         string  `yaml:"url"`
	Cron        *string `yaml:"cron"`
}

var Config Cfg

func LoadConfig() Cfg {
	file, err := os.ReadFile("sano.yml")
	if err != nil {
		log.Panicln("Cannot read [sano.yml] file.")
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Panicln("Invalid [sano.yml] structure,", err)
	}

	return Config
}

func CheckConfig() {
	_, err := os.ReadFile("sano.yml")
	if err != nil {
		log.Println("[sano.yml] file was not found, download it from GitHub.")
		utils.DownloadFile("https://raw.githubusercontent.com/hypnodev/sano-backend/main/sano.yml.default", "sano.yml")
		log.Panicln("[sano.yml] was created, please configure it before run Sano again.")
	}
}

func (cfg Cfg) GetService(name string) *Service {
	for _, service := range cfg.Services {
		if service.Name == name {
			return &service
		}
	}

	return nil
}
