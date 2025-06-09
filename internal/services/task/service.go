package task

import (
	"github.com/longgggwwww/hrm-ms-hr/ent"
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
	Client *ent.Client
}

// NewTaskService creates a new task service
func NewTaskService(client *ent.Client) *TaskService {
	return &TaskService{
		Client: client,
	}
}
