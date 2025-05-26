package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("code").
			NotEmpty().
			Annotations(entproto.Field(3)),
		field.Int("org_id").
			Annotations(entproto.Field(4)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(5)),
		field.Time("updated_at").
			Default(time.Now).
			Annotations(entproto.Field(6)),
	}
}

// Edges of the Department
func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("positions", Position.Type).
			Annotations(entproto.Field(7)),
		edge.From("organization", Organization.Type).
			Ref("departments").
			Field("org_id").
			Unique().
			Required().
			Annotations(entproto.Field(8)),
	}
}

func (Department) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code", "org_id").Unique(),
	}
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
