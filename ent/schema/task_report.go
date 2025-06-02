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
		field.String("content").
			Optional().
			Annotations(entproto.Field(3)),
		field.Int("task_id").
			Annotations(entproto.Field(4)),
		field.Int("reporter_id").
			Annotations(entproto.Field(5)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(6)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(7)),
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
			Annotations(entproto.Field(8)),
		edge.From("reporter", Employee.Type).
			Ref("task_reports").
			Field("reporter_id").
			Unique().
			Required().
			Annotations(entproto.Field(9)),
	}
}

func (TaskReport) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
