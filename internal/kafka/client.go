package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaClient wraps kafka writer for producing messages
type KafkaClient struct {
	writer *kafka.Writer
}

// NewKafkaClient creates a new Kafka client
func NewKafkaClient(brokers []string) *KafkaClient {
	if len(brokers) == 0 {
		log.Println("Warning: No Kafka brokers provided")
		return nil
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		BatchTimeout:           10 * time.Millisecond,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}

	return &KafkaClient{
		writer: writer,
	}
}

// PublishEvent publishes an event to the specified topic
func (c *KafkaClient) PublishEvent(ctx context.Context, topic string, key string, event interface{}) error {
	if c == nil || c.writer == nil {
		log.Printf("Kafka client not available, skipping event: %s", topic)
		return nil
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: eventData,
		Time:  time.Now(),
	}

	return c.writer.WriteMessages(ctx, message)
}

// Close closes the Kafka writer
func (c *KafkaClient) Close() error {
	if c != nil && c.writer != nil {
		return c.writer.Close()
	}
	return nil
}
