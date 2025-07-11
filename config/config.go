package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	Mail mailConfig
}

type mailConfig struct {
	URL                        string
	Port                       int
	Username, Password, Folder string
	PollingInterval            time.Duration
}

func NewConfig() *ServiceConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := new(ServiceConfig)

	// MAIL PART
	cfg.Mail.Username = os.Getenv("MAIL_USER")
	cfg.Mail.Password = os.Getenv("MAIL_PASS")
	cfg.Mail.Folder = os.Getenv("MAIL_FOLDER")
	cfg.Mail.URL = os.Getenv("MAIL_URL")

	portStr := os.Getenv("MAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Errorf("Config loader - Error loading port: %v", err))
	}
	cfg.Mail.Port = port

	pollingIntervalStr := os.Getenv("MAIL_POLLING_INTERVAL_SEC")
	pollingInterval, err := strconv.Atoi(pollingIntervalStr)
	if err != nil {
		panic(fmt.Errorf("Error during configuration loading: %v", err))
	}
	cfg.Mail.PollingInterval = time.Duration(pollingInterval) * time.Second

	// REDIS PART
	// TBD

	// KAFKA PART
	// TBD

	return cfg
}
