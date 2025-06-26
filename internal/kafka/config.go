package kafka

import (
	"os"
	"strings"
)

// Config holds Kafka configuration
type Config struct {
	Brokers []string
	Enabled bool
}

// NewConfig creates a new Kafka configuration from environment variables
func NewConfig() *Config {
	brokersStr := os.Getenv("KAFKA_BROKER")
	if brokersStr == "" {
		return &Config{
			Brokers: []string{},
			Enabled: false,
		}
	}

	brokers := strings.Split(brokersStr, ",")
	// Trim whitespace from each broker
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}

	return &Config{
		Brokers: brokers,
		Enabled: true,
	}
}

// Topics contains the Kafka topics used by the application
const (
	TopicTaskEvents       = "task-events"
	TopicTaskReportEvents = "task-report-events"
)
