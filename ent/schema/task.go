package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(entproto.Field(2)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(3)),
		field.Int("process").
			Default(0).
			Annotations(entproto.Field(4)),
		field.Enum("status").
			Values("not_started", "in_progress", "completed").
			Default("not_started").
			Annotations(entproto.Field(5), entproto.Enum(map[string]int32{
				"not_started": 0,
				"in_progress": 1,
				"completed":   2,
			})),
		field.Time("start_at").
			Annotations(entproto.Field(6)),
		field.Int("project_id").
			Annotations(entproto.Field(7)),
		field.Int("org_id").
			Annotations(entproto.Field(8)),
		field.Int("creator_id").
			Annotations(entproto.Field(9)),
		field.Int("updater_id").
			Annotations(entproto.Field(10)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(11)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(12)),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("tasks").
			Field("project_id").
			Unique().
			Required().
			Annotations(entproto.Field(13)),
	}
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
