package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type MailListenerServiceConfig struct {
	Mail  *MailConfig
	Kafka *KafkaConfig
	Redis *RedisConfig
}

type RunnerServiceConfig struct {
	Kafka *KafkaConfig
}

type MailConfig struct {
	Host                       string
	Username, Password, Folder string
	PollingInterval            time.Duration
}

type RedisConfig struct {
	Host           string
	User, Password string
	TTL            time.Duration
}

type KafkaConfig struct {
	Host    string
	Topic   string
	GroupID string
}

type envLoader struct {
	dotEnvPaths []string
}

func NewEnvLoader(paths []string) *envLoader {
	c := new(envLoader)
	c.dotEnvPaths = paths

	err := godotenv.Load(c.dotEnvPaths...)
	if err != nil {
		panic("Error loading .env file")
	}

	return c
}

func NewMailConfig(e *envLoader) *MailConfig {
	cfg := new(MailConfig)

	cfg.Host = os.Getenv("MAIL_HOST")
	cfg.Username = os.Getenv("MAIL_USER")
	cfg.Password = os.Getenv("MAIL_PASS")
	cfg.Folder = os.Getenv("MAIL_FOLDER")

	pollingIntervalStr := os.Getenv("MAIL_POLLING_INTERVAL_SEC")
	pollingInterval, err := strconv.Atoi(pollingIntervalStr)
	if err != nil {
		panic(fmt.Errorf("Error during configuration loading: %v", err))
	}
	cfg.PollingInterval = time.Duration(pollingInterval) * time.Second

	return cfg
}

func NewKafkaConfig(e *envLoader) *KafkaConfig {
	cfg := new(KafkaConfig)

	cfg.Host = os.Getenv("KAFKA_HOST")
	cfg.Topic = os.Getenv("KAFKA_TOPIC")
	cfg.GroupID = os.Getenv("KAFKA_CONSUMER_GROUP")

	return cfg
}

func NewRedisConfig(e *envLoader) *RedisConfig {
	cfg := new(RedisConfig)

	cfg.Host = os.Getenv("REDIS_HOST")
	cfg.User = os.Getenv("REDIS_USER")
	cfg.Password = os.Getenv("REDIS_PASS")
	ttlHoursStr := os.Getenv("REDIS_TTL_HOURS")
	ttlHours, err := strconv.Atoi(ttlHoursStr)
	if err != nil {
		panic(fmt.Errorf("Error during configuration loading: %v", err))
	}
	cfg.TTL = time.Duration(ttlHours) * time.Hour

	return cfg
}

func NewMailListenerServiceConfig(e *envLoader) *MailListenerServiceConfig {
	cfg := new(MailListenerServiceConfig)
	cfg.Kafka = NewKafkaConfig(e)
	cfg.Mail = NewMailConfig(e)
	cfg.Redis = NewRedisConfig(e)

	return cfg
}

func NewRunnerServiceConfig(e *envLoader) *RunnerServiceConfig {
	cfg := new(RunnerServiceConfig)
	cfg.Kafka = NewKafkaConfig(e)
	return cfg
}
