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
			StructTag(`json:"name"`).
			Annotations(entproto.Field(2)),
		field.String("code").
			Unique().
			StructTag(`json:"code"`).
			Annotations(entproto.Field(3)),
		field.String("description").
			Optional().
			StructTag(`json:"description"`).
			Annotations(entproto.Field(4)),
		field.Int("process").
			Default(0).
			StructTag(`json:"process"`).
			Annotations(entproto.Field(5)),
		field.Enum("status").
			Values("not_received", "received", "in_progress", "completed", "cancelled").
			Default("not_received").
			StructTag(`json:"status"`).
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
			Nillable().
			StructTag(`json:"start_at"`).
			Annotations(entproto.Field(7)),
		field.Time("due_date").
			Optional().
			Nillable().
			StructTag(`json:"due_date"`).
			Annotations(entproto.Field(8)),
		field.Int("project_id").
			Optional().
			StructTag(`json:"project_id"`).
			Annotations(entproto.Field(9)),
		field.Int("creator_id").
			StructTag(`json:"creator_id"`).
			Annotations(entproto.Field(10)),
		field.Int("updater_id").
			StructTag(`json:"updater_id"`).
			Annotations(entproto.Field(11)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(12)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(13)),
		field.Enum("type").
			Values("task", "feature", "bug", "another").
			Default("task").
			StructTag(`json:"type"`).
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
			StructTag(`json:"labels"`).
			Annotations(entproto.Field(16)), // Thêm edge tới label
		edge.To("assignees", Employee.Type).
			StructTag(`json:"assignees"`).
			Annotations(entproto.Field(17)), // Edge many-to-many với Employee
		edge.To("reports", TaskReport.Type).
			StructTag(`json:"reports"`).
			Annotations(entproto.Field(18)), // Edge đến TaskReport
	}
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
