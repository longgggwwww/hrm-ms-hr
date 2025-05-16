package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Organization holds the schema definition for the Organization entity.
type Organization struct {
	ent.Schema
}

// Fields of the Organization.
func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("code").
			NotEmpty().
			Unique().
			Annotations(entproto.Field(3)),
		field.Int("parent_id").
			Optional().
			Annotations(entproto.Field(4)),
		field.String("logo_url").
			Optional().
			Annotations(entproto.Field(5)),
		field.String("phone").
			Optional().
			Annotations(entproto.Field(6)),
		field.String("email").
			Optional().
			Annotations(entproto.Field(7)),
		field.String("website").
			Optional().
			Annotations(entproto.Field(8)),
		field.String("address").
			Optional().
			Annotations(entproto.Field(9)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(10)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(11)),
	}
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).
			Ref("organizations").
			Unique().
			Field("parent_id").
			Annotations(entproto.Field(8)),
		edge.To("organizations", Organization.Type).
			Annotations(entproto.Field(9)),
	}
}

func (Organization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
