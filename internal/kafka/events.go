package kafka

import (
	"time"

	"github.com/longgggwwww/hrm-ms-hr/ent"
)

// TaskEventType represents the type of task event
type TaskEventType string

const (
	TaskCreated TaskEventType = "task.created"
	TaskUpdated TaskEventType = "task.updated"
)

// TaskReportEventType represents the type of task report event
type TaskReportEventType string

const (
	TaskReportCreated TaskReportEventType = "task_report.created"
	TaskReportUpdated TaskReportEventType = "task_report.updated"
)

// TaskEvent represents a task event
type TaskEvent struct {
	EventID      string        `json:"event_id"`
	EventType    TaskEventType `json:"event_type"`
	Timestamp    time.Time     `json:"timestamp"`
	Source       string        `json:"source"`
	TaskID       int           `json:"task_id"`
	TaskCode     string        `json:"task_code"`
	TaskName     string        `json:"task_name"`
	Description  string        `json:"description,omitempty"`
	ProjectID    *int          `json:"project_id,omitempty"`
	DepartmentID *int          `json:"department_id,omitempty"`
	CreatorID    int           `json:"creator_id"`
	UpdaterID    int           `json:"updater_id"`
	Status       string        `json:"status"`
	Type         string        `json:"type"`
	Process      int           `json:"process"`
	StartAt      *time.Time    `json:"start_at,omitempty"`
	DueDate      *time.Time    `json:"due_date,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	AssigneeIDs  []int         `json:"assignee_ids,omitempty"`
	LabelIDs     []int         `json:"label_ids,omitempty"`
	OrgID        int           `json:"org_id,omitempty"`
	ZaloGid      *string       `json:"zalo_gid,omitempty"`
}

// TaskReportEvent represents a task report event
type TaskReportEvent struct {
	EventID    string              `json:"event_id"`
	EventType  TaskReportEventType `json:"event_type"`
	Timestamp  time.Time           `json:"timestamp"`
	Source     string              `json:"source"`
	ReportID   int                 `json:"report_id"`
	TaskID     int                 `json:"task_id"`
	ReporterID int                 `json:"reporter_id"`
	Content    string              `json:"content"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	OrgID      int                 `json:"org_id,omitempty"`
}

// NewTaskCreatedEvent creates a new task created event
func NewTaskCreatedEvent(task *ent.Task, orgID int) *TaskEvent {
	assigneeIDs := make([]int, 0)
	if task.Edges.Assignees != nil {
		for _, assignee := range task.Edges.Assignees {
			assigneeIDs = append(assigneeIDs, assignee.ID)
		}
	}

	labelIDs := make([]int, 0)
	if task.Edges.Labels != nil {
		for _, label := range task.Edges.Labels {
			labelIDs = append(labelIDs, label.ID)
		}
	}

	// Convert int to *int for ProjectID
	var projectID *int
	if task.ProjectID != 0 {
		projectID = &task.ProjectID
	}

	// Convert int to *int for DepartmentID
	var departmentID *int
	if task.DepartmentID != 0 {
		departmentID = &task.DepartmentID
	}

	// Get zalo_gid from Department edge
	var zaloGid *string
	if task.Edges.Department != nil && task.Edges.Department.ZaloGid != nil {
		zaloGid = task.Edges.Department.ZaloGid
	}

	return &TaskEvent{
		EventID:      generateEventID(),
		EventType:    TaskCreated,
		Timestamp:    time.Now(),
		Source:       "hrm-ms-hr",
		TaskID:       task.ID,
		TaskCode:     task.Code,
		TaskName:     task.Name,
		Description:  task.Description,
		ProjectID:    projectID,
		DepartmentID: departmentID,
		CreatorID:    task.CreatorID,
		UpdaterID:    task.UpdaterID,
		Status:       string(task.Status),
		Type:         string(task.Type),
		Process:      task.Process,
		StartAt:      task.StartAt,
		DueDate:      task.DueDate,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
		AssigneeIDs:  assigneeIDs,
		LabelIDs:     labelIDs,
		OrgID:        orgID,
		ZaloGid:      zaloGid,
	}
}

// NewTaskUpdatedEvent creates a new task updated event
func NewTaskUpdatedEvent(task *ent.Task, orgID int) *TaskEvent {
	assigneeIDs := make([]int, 0)
	if task.Edges.Assignees != nil {
		for _, assignee := range task.Edges.Assignees {
			assigneeIDs = append(assigneeIDs, assignee.ID)
		}
	}

	labelIDs := make([]int, 0)
	if task.Edges.Labels != nil {
		for _, label := range task.Edges.Labels {
			labelIDs = append(labelIDs, label.ID)
		}
	}

	// Convert int to *int for ProjectID
	var projectID *int
	if task.ProjectID != 0 {
		projectID = &task.ProjectID
	}

	// Convert int to *int for DepartmentID
	var departmentID *int
	if task.DepartmentID != 0 {
		departmentID = &task.DepartmentID
	}

	// Get zalo_gid from Department edge
	var zaloGid *string
	if task.Edges.Department != nil && task.Edges.Department.ZaloGid != nil {
		zaloGid = task.Edges.Department.ZaloGid
	}

	return &TaskEvent{
		EventID:      generateEventID(),
		EventType:    TaskUpdated,
		Timestamp:    time.Now(),
		Source:       "hrm-ms-hr",
		TaskID:       task.ID,
		TaskCode:     task.Code,
		TaskName:     task.Name,
		Description:  task.Description,
		ProjectID:    projectID,
		DepartmentID: departmentID,
		CreatorID:    task.CreatorID,
		UpdaterID:    task.UpdaterID,
		Status:       string(task.Status),
		Type:         string(task.Type),
		Process:      task.Process,
		StartAt:      task.StartAt,
		DueDate:      task.DueDate,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
		AssigneeIDs:  assigneeIDs,
		LabelIDs:     labelIDs,
		OrgID:        orgID,
		ZaloGid:      zaloGid,
	}
}

// NewTaskReportCreatedEvent creates a new task report created event
func NewTaskReportCreatedEvent(report *ent.TaskReport, orgID int) *TaskReportEvent {
	return &TaskReportEvent{
		EventID:    generateEventID(),
		EventType:  TaskReportCreated,
		Timestamp:  time.Now(),
		Source:     "hrm-ms-hr",
		ReportID:   report.ID,
		TaskID:     report.TaskID,
		ReporterID: report.ReporterID,
		Content:    report.Content,
		CreatedAt:  report.CreatedAt,
		UpdatedAt:  report.UpdatedAt,
		OrgID:      orgID,
	}
}

// NewTaskReportUpdatedEvent creates a new task report updated event
func NewTaskReportUpdatedEvent(report *ent.TaskReport, orgID int) *TaskReportEvent {
	return &TaskReportEvent{
		EventID:    generateEventID(),
		EventType:  TaskReportUpdated,
		Timestamp:  time.Now(),
		Source:     "hrm-ms-hr",
		ReportID:   report.ID,
		TaskID:     report.TaskID,
		ReporterID: report.ReporterID,
		Content:    report.Content,
		CreatedAt:  report.CreatedAt,
		UpdatedAt:  report.UpdatedAt,
		OrgID:      orgID,
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
