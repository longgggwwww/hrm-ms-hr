package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TaskReport holds the schema definition for the TaskReport entity.
type TaskReport struct {
	ent.Schema
}

// Fields of the TaskReport.
func (TaskReport) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			Annotations(entproto.Field(2)),
		field.String("content").
			Optional().
			Annotations(entproto.Field(3)),
		field.Enum("status").
			Values("not_received", "received", "in_progress", "completed", "cancelled").
			Annotations(
				entproto.Field(4),
				entproto.Enum(map[string]int32{
					"not_received": 0,
					"received":     1,
					"in_progress":  2,
					"completed":    3,
					"cancelled":    4,
				})),
		field.Int("progress_percentage").
			Default(0).
			Min(0).
			Max(100).
			Annotations(entproto.Field(5)),
		field.Time("reported_at").
			Default(time.Now).
			Annotations(entproto.Field(6)),
		field.String("issues_encountered").
			Optional().
			Annotations(entproto.Field(7)),
		field.String("next_steps").
			Optional().
			Annotations(entproto.Field(8)),
		field.Time("estimated_completion").
			Optional().
			Annotations(entproto.Field(9)),
		field.Int("task_id").
			Annotations(entproto.Field(10)),
		field.Int("reporter_id").
			Annotations(entproto.Field(11)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(12)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(13)),
	}
}

// Edges of the TaskReport.
func (TaskReport) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).
			Ref("reports").
			Field("task_id").
			Unique().
			Required().
			Annotations(entproto.Field(14)),
		edge.From("reporter", Employee.Type).
			Ref("task_reports").
			Field("reporter_id").
			Unique().
			Required().
			Annotations(entproto.Field(15)),
	}
}

func (TaskReport) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
