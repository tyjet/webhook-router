package event

import (
	log "github.com/tyjet/soos-ult"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Producer interface {
	Produce(topic, key string, value []byte) error
	HandleEvents()
}

type producer struct {
	*kafka.Producer
}

func NewProducer(p *kafka.Producer) Producer {
	return &producer{Producer: p}
}

func (p *producer) Produce(topic, key string, value []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}
	return p.Producer.Produce(msg, nil)
}

func (p *producer) HandleEvents() {
	for e := range p.Events() {
		switch event := e.(type) {
		case *kafka.Message:
			if event.TopicPartition.Error != nil {
				log.Error("failed to produce message", zap.Error(event.TopicPartition.Error), zap.String("topc", *event.TopicPartition.Topic), zap.Int32("partition", event.TopicPartition.Partition))
			} else {
				log.Debug("successfully produced message", zap.String("topic", *event.TopicPartition.Topic), zap.Int32("partition", event.TopicPartition.Partition), zap.Int64("offset", int64(event.TopicPartition.Offset)))
			}
		}
	}
}
