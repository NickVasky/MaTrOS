package client

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/NickVasky/MaTrOS/shared/config"
	"github.com/joho/godotenv"
)

func loadEnvs() (*config.OrchConfig, error) {
	cfg := new(config.OrchConfig)
	envs := make(map[string]string)
	err := godotenv.Load("test.env")
	if err != nil {
		return cfg, err
	}

	envs = map[string]string{"TEST_HOST": "", "TEST_USER": "", "TEST_PASS": ""}
	for k, _ := range envs {
		envValue := os.Getenv(k)
		if envValue == "" {
			return cfg, fmt.Errorf("Env '%v' is empty", k)
		}
		envs[k] = envValue
	}

	cfg.Host = envs["TEST_HOST"]
	cfg.Password = envs["TEST_PASS"]
	cfg.User = envs["TEST_USER"]
	cfg.URLSchema = "https"
	cfg.RobotEdition = 2

	return cfg, nil
}

func TestAccount(t *testing.T) {
	cfg, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(cfg)

	if err != nil {
		t.Errorf("PostAccount() err: %v", err)
	}
	log.Println("Got token:", bot.apiToken)

}

func TestGetProjects(t *testing.T) {
	cfg, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(cfg)

	if err != nil {
		t.Errorf("PostAccount() err: %v", err)
	}
	projects, err := bot.GetProjects()
	if err != nil {
		t.Errorf("GetProjects() err: %v", err)
	}
	log.Println(projects)
}

func TestGetRobots(t *testing.T) {
	cfg, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(cfg)

	if err != nil {
		t.Errorf("PostAccount() err: %v", err)
	}
	robots, err := bot.GetRobots()
	if err != nil {
		t.Errorf("GetRobots() err: %v", err)
	}
	log.Println(robots)
}

func TestPutRobotStartAsync(t *testing.T) {
	cfg, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(cfg)

	if err != nil {
		t.Errorf("PostAccount() err: %v", err)
	}
	err = bot.PutRobotStartAsync(1, 1)
	if err != nil {
		t.Errorf("PutRobotStartAsync() err: %v", err)
	}
}
