package kafka

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// Consumer is a Kafka consumer
type Consumer struct {
	reader *kafka.Reader
}

// Producer is a Kafka producer
type Producer struct {
	writer *kafka.Writer
}

// NewConsumer creates a new Kafka consumer for a topic
func NewConsumer(topic string) (*Consumer, error) {
	// Get Kafka brokers from environment variable
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	// Get Kafka username and password from environment variables
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	// Set up the reader configuration
	config := kafka.ReaderConfig{
		Brokers:        strings.Split(brokers, ","),
		GroupID:        fmt.Sprintf("bike-parts-finder-%s", topic),
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
	}

	// If authentication is required
	if username != "" && password != "" {
		mechanism := plain.Mechanism{
			Username: username,
			Password: password,
		}
		dialer := &kafka.Dialer{
			Timeout:       10 * time.Second,
			SASLMechanism: mechanism,
		}
		config.Dialer = dialer
	}

	// Create the reader
	reader := kafka.NewReader(config)

	return &Consumer{
		reader: reader,
	}, nil
}

// ReadMessage reads a message from Kafka
func (c *Consumer) ReadMessage(timeout time.Duration) (kafka.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.reader.ReadMessage(ctx)
}

// CommitMessage commits the offset of a message
func (c *Consumer) CommitMessage(msg kafka.Message) error {
	return c.reader.CommitMessages(context.Background(), msg)
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// NewProducer creates a new Kafka producer for a topic
func NewProducer(topic string) (*Producer, error) {
	// Get Kafka brokers from environment variable
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	// Get Kafka username and password from environment variables
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	// Set up the writer configuration
	var dialer *kafka.Dialer
	if username != "" && password != "" {
		mechanism := plain.Mechanism{
			Username: username,
			Password: password,
		}
		dialer = &kafka.Dialer{
			Timeout:       10 * time.Second,
			SASLMechanism: mechanism,
		}
	}

	// Create the writer
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(brokers, ",")...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           10 * time.Second,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}

	if dialer != nil {
		writer.Transport = &kafka.Transport{
			SASL: dialer.SASLMechanism,
		}
	}

	return &Producer{
		writer: writer,
	}, nil
}

// WriteMessage writes a message to Kafka
func (p *Producer) WriteMessage(value []byte) error {
	return p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: value,
		},
	)
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}
