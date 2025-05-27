package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(entproto.Field(2)),
		field.String("code").
			Annotations(entproto.Field(3)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(4)),
		field.Time("start_at").
			Annotations(entproto.Field(5)),
		field.Time("end_at").
			Optional().
			Annotations(entproto.Field(6)),
		field.Int("creator_id").
			Annotations(entproto.Field(7)),
		field.Int("updater_id").
			Annotations(entproto.Field(8)),
		field.Int("org_id").
			Annotations(entproto.Field(9)),
		field.Int("process").
			Optional().
			Annotations(entproto.Field(10)),
		field.Enum("status").
			Values("not_started", "in_progress", "completed").
			Default("not_started").
			Annotations(
				entproto.Field(11),
				entproto.Enum(map[string]int32{
					"not_started": 0,
					"in_progress": 1,
					"completed":   2,
				})),
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

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Annotations(entproto.Field(14)),
		edge.From("organization", Organization.Type).
			Ref("projects").
			Field("org_id").
			Unique().
			Required().
			Annotations(entproto.Field(15)),
		edge.From("creator", Employee.Type).
			Ref("created_projects").
			Field("creator_id").
			Unique().
			Required().
			Annotations(entproto.Field(16)),
		edge.From("updater", Employee.Type).
			Ref("updated_projects").
			Field("updater_id").
			Unique().
			Required().
			Annotations(entproto.Field(17)),
	}
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
