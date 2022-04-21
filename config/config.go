package config

import (
	"errors"
	"github.com/Shopify/sarama"
	"os"
	"strconv"
	"strings"
)

type ProducerConfig struct {
	BootstrapServer []string
	Topics []string
	ProduceBackoffMs int
	Config *sarama.Config
}

func LoadConfigOrDie() (*ProducerConfig, error) {
	c := sarama.NewConfig()
	bootstrap := os.Getenv("BOOTSTRAP_SERVER")
	if bootstrap == "" {
		return nil, errors.New("BOOTSTRAP_SERVER must be set")
	}

	bootstrapList := strings.Split(bootstrap, ",")
	topics := os.Getenv("TOPICS")
	if topics == "" {
		return nil, errors.New("TOPICS must be set")
	}

	topicList := strings.Split(topics, ",")
	backoff := os.Getenv("PRODUCE_BACKOFF_MS")
	var (
		err error
		backoffMs = 1000
	)
	if backoff != "" {
		backoffMs, err = strconv.Atoi(backoff)
		if err != nil {
			return nil, err
		}
	}

	return &ProducerConfig{
		BootstrapServer:  bootstrapList,
		Topics:           topicList,
		ProduceBackoffMs: backoffMs,
		Config:           c,
	}, nil
}