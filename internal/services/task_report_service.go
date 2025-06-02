package services

import (
	"context"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/taskreport"
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
	return client.TaskReport.Query().
		Where(taskreport.ID(report.ID)).
		WithTask().
		WithReporter().
		Only(ctx)
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
	return client.TaskReport.Query().
		Where(taskreport.ID(id)).
		WithTask().
		WithReporter().
		Only(ctx)
}
