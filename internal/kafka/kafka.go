package kafka

import (
	log "github.com/tyjet/soos-ult"
	"github.com/tyjet/webhook-router/internal/config"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func Setup() (*kafka.Producer, bool) {
	bootstrapServers := config.KafkaBootstrapServers()
	return NewProducer(bootstrapServers)
}

func NewProducer(bootstrapServers string) (*kafka.Producer, bool) {
	cfgMap := &kafka.ConfigMap{"bootstrap.servers": bootstrapServers}
	p, err := kafka.NewProducer(cfgMap)
	if err != nil {
		log.Error("could not create kafka producer", zap.Error(err))
		return nil, false
	}

	return p, true
}
