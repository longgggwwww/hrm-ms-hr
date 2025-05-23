package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("title").
			Annotations(entproto.Field(2)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(3)),
		field.Int("process").
			Annotations(entproto.Field(4)),
		field.Bool("status").
			Annotations(entproto.Field(5)),
		field.Time("start_at").
			Annotations(entproto.Field(6)),
		field.UUID("project_id", uuid.UUID{}).
			Annotations(entproto.Field(7)),
		field.UUID("branch_id", uuid.NullUUID{}).
			Annotations(entproto.Field(8)),
		field.UUID("creator_id", uuid.UUID{}).
			Annotations(entproto.Field(9)),
		field.UUID("updater_id", uuid.NullUUID{}).
			Annotations(entproto.Field(10)),
		field.Time("created_at").
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
		edge.To("labels", Label.Type), // Thêm edge tới label
	}
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
