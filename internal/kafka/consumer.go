package kafka

import (
	"context"
	"encoding/json"
	"firstmod/internal/models"
	"firstmod/internal/ports"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerImpl struct {
	reader  *kafka.Reader
	service ports.OrderService
	log     *slog.Logger
}

func NewConsumer(log *slog.Logger, brokers []string, topic, groupID string, service ports.OrderService) *KafkaConsumerImpl {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		Logger:         kafka.LoggerFunc(func(msg string, args ...interface{}) { log.Debug(msg, args...) }),
		ErrorLogger:    kafka.LoggerFunc(func(msg string, args ...interface{}) { log.Error(msg, args...) }),
	})
	log.Info("Kafka consumer initialized", "brokers", brokers, "topic", topic, "group_id", groupID)
	return &KafkaConsumerImpl{reader: reader, service: service, log: log}
}

func (c *KafkaConsumerImpl) StartConsuming(ctx context.Context) {
	c.log.Info("starting Kafka consumer")
	for {
		select {
		case <-ctx.Done():
			c.log.Info("Kafka consumer shutting down")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.log.Error("failed to fetch message from Kafka", "error", err)
				if ctx.Err() != nil {
					return
				}
				time.Sleep(time.Second)
				continue
			}

			c.log.Debug("received message from Kafka", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key))

			var order models.Order
			err = json.Unmarshal(msg.Value, &order)
			if err != nil {
				c.log.Error("failed to unmarshal Kafka message value to Order model", "offset", msg.Offset, "error", err, "value", string(msg.Value))
				if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
					c.log.Error("failed to commit invalid message", "offset", msg.Offset, "error", commitErr)
				}
				continue
			}
			err = c.service.Add(ctx, order)
			if err != nil {
				c.log.Error("failed to add order from Kafka message via service", "order_uid", order.OrderUID, "offset", msg.Offset, "error", err)
				if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
					c.log.Error("failed to commit message after processing error", "offset", msg.Offset, "error", commitErr)
				}
				continue
			}

			if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
				c.log.Error("failed to commit message after successful processing", "offset", msg.Offset, "error", commitErr)
			}
			c.log.Info("order processed and committed from Kafka", "order_uid", order.OrderUID, "offset", msg.Offset)
		}
	}
}

func (c *KafkaConsumerImpl) Close() error {
	c.log.Info("closing Kafka consumer")
	return c.reader.Close()
}
