package services

import (
	"context"
	"log"
	"strconv"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/taskreport"
	"github.com/longgggwwww/hrm-ms-hr/internal/kafka"
)

func CreateTaskReport(ctx context.Context, client *ent.Client, taskID, reporterID int, content string) (*ent.TaskReport, error) {
	report, err := client.TaskReport.Create().
		SetTaskID(taskID).
		SetReporterID(reporterID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	result, err := client.TaskReport.Query().
		Where(taskreport.ID(report.ID)).
		WithTask().
		WithReporter().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Publish task report created event
	publishTaskReportCreatedEvent(ctx, result)

	return result, nil
}

func UpdateTaskReport(ctx context.Context, client *ent.Client, id, reporterID int, content string) (*ent.TaskReport, error) {
	if _, err := client.TaskReport.Query().Where(taskreport.ID(id), taskreport.ReporterID(reporterID)).Only(ctx); err != nil {
		return nil, err
	}
	_, err := client.TaskReport.UpdateOneID(id).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	result, err := client.TaskReport.Query().
		Where(taskreport.ID(id)).
		WithTask().
		WithReporter().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Publish task report updated event
	publishTaskReportUpdatedEvent(ctx, result)

	return result, nil
}

// Global Kafka client variable
var kafkaClient *kafka.KafkaClient

// SetKafkaClient sets the global Kafka client for task report services
func SetTaskReportKafkaClient(client *kafka.KafkaClient) {
	kafkaClient = client
}

// publishTaskReportCreatedEvent publishes a task report created event to Kafka
func publishTaskReportCreatedEvent(ctx context.Context, report *ent.TaskReport) {
	if kafkaClient == nil {
		return
	}

	// Get organization ID from reporter
	orgID := 0
	if report.Edges.Reporter != nil {
		orgID = report.Edges.Reporter.OrgID
	}

	event := kafka.NewTaskReportCreatedEvent(report, orgID)
	key := strconv.Itoa(report.ID)

	if err := kafkaClient.PublishEvent(ctx, kafka.TopicTaskReportEvents, key, event); err != nil {
		log.Printf("Failed to publish task report created event: %v", err)
	}
}

// publishTaskReportUpdatedEvent publishes a task report updated event to Kafka
func publishTaskReportUpdatedEvent(ctx context.Context, report *ent.TaskReport) {
	if kafkaClient == nil {
		return
	}

	// Get organization ID from reporter
	orgID := 0
	if report.Edges.Reporter != nil {
		orgID = report.Edges.Reporter.OrgID
	}

	event := kafka.NewTaskReportUpdatedEvent(report, orgID)
	key := strconv.Itoa(report.ID)

	if err := kafkaClient.PublishEvent(ctx, kafka.TopicTaskReportEvents, key, event); err != nil {
		log.Printf("Failed to publish task report updated event: %v", err)
	}
}
