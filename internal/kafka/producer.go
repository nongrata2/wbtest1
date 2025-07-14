package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaProducerImpl struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func NewProducer(log *slog.Logger, brokers []string, topic string) *KafkaProducerImpl {
	writer := &kafka.Writer{
		Addr:        kafka.TCP(brokers...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		Logger:      kafka.LoggerFunc(func(msg string, args ...interface{}) { log.Debug(msg, args...) }),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) { log.Error(msg, args...) }),
	}
	log.Info("Kafka producer initialized", "brokers", brokers, "topic", topic)
	return &KafkaProducerImpl{writer: writer, log: log}
}

func (p *KafkaProducerImpl) Publish(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Error("failed to publish message to Kafka", "key", key, "error", err)
		return err
	}
	p.log.Debug("message published to Kafka", "key", key)
	return nil
}

func (p *KafkaProducerImpl) Close() error {
	p.log.Info("closing Kafka producer")
	return p.writer.Close()
}
