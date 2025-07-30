package client

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestNewBotApiClient(t *testing.T) {
	type args struct {
		Host      string
		URLSchema string
	}
	tests := []struct {
		name string
		args args
		want *BotApiClient
	}{
		{
			name: "Test 1",
			args: args{Host: "localhost", URLSchema: "https"},
			want: &BotApiClient{
				Client:    &http.Client{Timeout: time.Second * 10, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
				URLSchema: "https",
				Host:      "localhost",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBotApiClient(tt.args.Host, tt.args.URLSchema); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBotApiClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount(t *testing.T) {
	err := godotenv.Load("test.env")
	if err != nil {
		t.Error(err)
	}

	envs := map[string]string{"TEST_HOST": "", "TEST_USER": "", "TEST_PASS": ""}
	for k, _ := range envs {
		envValue := os.Getenv(k)
		if envValue == "" {
			t.Errorf("Env '%v' is empty", k)
		}
		envs[k] = envValue
	}

	bot := NewBotApiClient(envs["TEST_HOST"], "https")
	t.Run("Test1", func(t *testing.T) {
		token, err := bot.PostAccount(envs["TEST_USER"], envs["TEST_PASS"], 2)
		if err != nil {
			t.Errorf("PostAccount() err: %v", err)
		}
		log.Println("Got token:", token)
	})
}
