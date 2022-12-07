package main

import (
	"encoding/json"
	"os"

	log "github.com/tyjet/soos-ult"
	"github.com/tyjet/webhook-router/internal/config"
	"github.com/tyjet/webhook-router/internal/event"
	"github.com/tyjet/webhook-router/internal/kafka"
	"go.uber.org/zap"
)

func main() {
	if err := log.Setup(); err != nil {
		os.Exit(1)
	}
	defer log.Sync()

	if ok := config.ReadConfig("config", config.YAML, []string{"etc"}); !ok {
		log.Sync()
		os.Exit(1)
	}

	log.Info("Parsed config", zap.String("kafka.bootstrap.servers", config.KafkaBootstrapServers()))
	kafkaProducer, ok := kafka.Setup()
	if !ok {
		log.Sync()
		os.Exit(1)
	}
	defer func() {
		kafkaProducer.Close()
		unflushed := kafkaProducer.Flush(1000)
		if unflushed != 0 {
			log.Error("could not flush all events", zap.Int("unflushed", unflushed))
		}
	}()

	producer := event.NewProducer(kafkaProducer)
	value, err := json.Marshal("ddiner-test")
	if err != nil {
		log.Error("could not marshal test message", zap.Error(err))
		log.Sync()
		kafkaProducer.Close()
		os.Exit(1)
	}

	err = producer.Produce("tests", "name", value)
	if err != nil {
		log.Error("could not produce test message", zap.Error(err))
	}
}
