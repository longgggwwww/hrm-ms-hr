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
		field.String("code").
			Unique().
			Annotations(entproto.Field(3)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(4)),
		field.Int("process").
			Default(0).
			Annotations(entproto.Field(5)),
		field.Enum("status").
			Values("not_received", "received", "in_progress", "completed", "cancelled").
			Default("not_received").
			Annotations(
				entproto.Field(6),
				entproto.Enum(map[string]int32{
					"not_received": 0,
					"received":     1,
					"in_progress":  2,
					"completed":    3,
					"cancelled":    4,
				})),
		field.Time("start_at").
			Optional().
			Annotations(entproto.Field(7)),
		field.Time("due_date").
			Optional().
			Annotations(entproto.Field(8)),
		field.Int("project_id").
			Optional().
			Annotations(entproto.Field(9)),
		field.Int("creator_id").
			Annotations(entproto.Field(10)),
		field.Int("updater_id").
			Annotations(entproto.Field(11)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(12)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(13)),
		field.Enum("type").
			Values("task", "feature", "bug", "another").
			Default("task").
			Annotations(
				entproto.Field(14),
				entproto.Enum(map[string]int32{
					"task":    0,
					"feature": 1,
					"bug":     2,
					"another": 3,
				})),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("tasks").
			Field("project_id").
			Unique().
			Annotations(entproto.Field(15)),
		edge.To("labels", Label.Type).
			Annotations(entproto.Field(16)), // Thêm edge tới label
		edge.To("assignees", Employee.Type).
			Annotations(entproto.Field(17)), // Edge many-to-many với Employee
	}
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
