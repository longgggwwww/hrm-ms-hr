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
			Annotations(entproto.Field(3)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(4)),
		field.Int("process").
			Annotations(entproto.Field(5)),
		field.Bool("status").
			Annotations(entproto.Field(6)),
		field.Time("start_at").
			Annotations(entproto.Field(7)),
		field.Int("project_id").
			Optional().
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
		field.Enum("type").
			Values("task", "feature", "bug", "another").
			Default("task").
			Annotations(
				entproto.Field(13),
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
			Annotations(entproto.Field(14)),
		edge.To("labels", Label.Type).
			Annotations(entproto.Field(15)), // Thêm edge tới label
	}
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
