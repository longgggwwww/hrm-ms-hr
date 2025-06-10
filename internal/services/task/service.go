package task

import (
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/kafka"
)

// ServiceError represents a service-level error
type ServiceError struct {
	Status int
	Msg    string
}

func (e *ServiceError) Error() string {
	return e.Msg
}

// TaskService provides business logic for task operations
type TaskService struct {
	Client      *ent.Client
	KafkaClient *kafka.KafkaClient
}

// NewTaskService creates a new task service
func NewTaskService(client *ent.Client) *TaskService {
	return &TaskService{
		Client:      client,
		KafkaClient: nil, // Will be set later
	}
}

// SetKafkaClient sets the Kafka client for the service
func (s *TaskService) SetKafkaClient(kafkaClient *kafka.KafkaClient) {
	s.KafkaClient = kafkaClient
}
