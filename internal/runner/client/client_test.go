package client

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func loadEnvs() (map[string]string, error) {
	envs := make(map[string]string)
	err := godotenv.Load("test.env")
	if err != nil {
		return envs, err
	}

	envs = map[string]string{"TEST_HOST": "", "TEST_USER": "", "TEST_PASS": ""}
	for k, _ := range envs {
		envValue := os.Getenv(k)
		if envValue == "" {
			return envs, fmt.Errorf("Env '%v' is empty", k)
		}
		envs[k] = envValue
	}
	return envs, nil
}

func TestAccount(t *testing.T) {
	envs, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(envs["TEST_HOST"], "https", envs["TEST_USER"], envs["TEST_PASS"], Standard)

	if err != nil {
		t.Errorf("PostAccount() err: %v", err)
	}
	log.Println("Got token:", bot.apiToken)

}

func TestGetProjects(t *testing.T) {
	envs, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(envs["TEST_HOST"], "https", envs["TEST_USER"], envs["TEST_PASS"], Standard)

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
	envs, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(envs["TEST_HOST"], "https", envs["TEST_USER"], envs["TEST_PASS"], Standard)

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
	envs, err := loadEnvs()
	if err != nil {
		t.Error(err)
	}

	bot, err := NewBotApiClient(envs["TEST_HOST"], "https", envs["TEST_USER"], envs["TEST_PASS"], Standard)

	if err != nil {
		t.Errorf("PostAccount() err: %v", err)
	}
	err = bot.PutRobotStartAsync(1, 1)
	if err != nil {
		t.Errorf("PutRobotStartAsync() err: %v", err)
	}
}
