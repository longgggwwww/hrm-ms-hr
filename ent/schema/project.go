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

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("name").
			Annotations(entproto.Field(2)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(3)),
		field.Time("start_at").
			Annotations(entproto.Field(4)),
		field.Time("end_at").
			Optional().
			Annotations(entproto.Field(5)),
		field.UUID("creator_id", uuid.UUID{}).
			Annotations(entproto.Field(6)),
		field.UUID("updater_id", uuid.NullUUID{}).
			Annotations(entproto.Field(7)),
		field.Time("created_at").
			Default(time.Now).
			Annotations(entproto.Field(8)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(9)),
		field.UUID("branch_id", uuid.NullUUID{}).
			Annotations(entproto.Field(10)),
		field.Int("process").
			Annotations(entproto.Field(11)),
		field.Enum("status").
			Values("not_started", "in_progress", "completed").
			Default("not_started").
			Annotations(
				entproto.Field(12),
				entproto.Enum(map[string]int32{
					"not_started": 0,
					"in_progress": 1,
					"completed":   2,
				})),
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Annotations(entproto.Field(13)),
		edge.To("labels", Label.Type), // Thêm edge tới label
	}
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
