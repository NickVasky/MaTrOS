package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type MailListenerServiceConfig struct {
	Mail  *MailConfig
	Kafka *KafkaConfig
	Redis *RedisConfig
}

type RunnerServiceConfig struct {
	Kafka *KafkaConfig
	Orch  *OrchConfig
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

type OrchConfig struct {
	URLSchema    string
	Host         string
	User         string
	Password     string
	RobotEdition uint
}

func NewMailConfig() *MailConfig {
	cfg := new(MailConfig)

	cfg.Host = os.Getenv("MAIL_HOST")
	cfg.Username = os.Getenv("MAIL_USER")
	cfg.Password = os.Getenv("MAIL_PASS")
	cfg.Folder = os.Getenv("MAIL_FOLDER")

	pollingIntervalStr := os.Getenv("MAIL_POLLING_INTERVAL_SEC")
	pollingInterval, err := strconv.Atoi(pollingIntervalStr)
	if err != nil {
		panic(fmt.Errorf("Error during Mail cfg loading: %v", err))
	}
	cfg.PollingInterval = time.Duration(pollingInterval) * time.Second

	return cfg
}

func NewKafkaConfig() *KafkaConfig {
	cfg := new(KafkaConfig)

	cfg.Host = os.Getenv("KAFKA_HOST")
	cfg.Topic = os.Getenv("KAFKA_TOPIC")
	cfg.GroupID = os.Getenv("KAFKA_CONSUMER_GROUP")

	return cfg
}

func NewRedisConfig() *RedisConfig {
	cfg := new(RedisConfig)

	cfg.Host = os.Getenv("REDIS_HOST")
	cfg.User = os.Getenv("REDIS_USER")
	cfg.Password = os.Getenv("REDIS_PASS")
	ttlHoursStr := os.Getenv("REDIS_TTL_HOURS")
	ttlHours, err := strconv.Atoi(ttlHoursStr)
	if err != nil {
		panic(fmt.Errorf("Error during Redis cfg loading: %v", err))
	}
	cfg.TTL = time.Duration(ttlHours) * time.Hour

	return cfg
}

func NewOrchConfig() *OrchConfig {
	cfg := new(OrchConfig)

	hostString := os.Getenv("ORCH_HOST")
	hostUrl, err := url.Parse(hostString)
	if err != nil {
		panic(fmt.Errorf("Error during Orch cfg loading: %v", err))
	}
	cfg.Host = hostUrl.Host
	cfg.URLSchema = hostUrl.Scheme

	cfg.User = os.Getenv("ORCH_USER")
	cfg.Password = os.Getenv("ORCH_PASS")

	robotEditionString := os.Getenv("ORCH_ROBOT_EDITION")
	var robotEdition uint64
	if robotEditionString == "" {
		robotEdition = 2 //default
	} else {
		robotEdition, err = strconv.ParseUint(robotEditionString, 10, 0)
		if err != nil {
			panic(fmt.Errorf("Error during Orch cfg loading: %v", err))
		}
	}
	cfg.RobotEdition = uint(robotEdition)

	return cfg
}

func NewMailListenerServiceConfig() *MailListenerServiceConfig {
	cfg := new(MailListenerServiceConfig)
	cfg.Kafka = NewKafkaConfig()
	cfg.Mail = NewMailConfig()
	cfg.Redis = NewRedisConfig()

	return cfg
}

func NewRunnerServiceConfig() *RunnerServiceConfig {
	cfg := new(RunnerServiceConfig)
	cfg.Kafka = NewKafkaConfig()
	cfg.Orch = NewOrchConfig()
	return cfg
}
