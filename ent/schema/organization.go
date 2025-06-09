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
			StructTag(`json:"name"`).
			Annotations(entproto.Field(2)),
		field.String("code").
			NotEmpty().
			Unique().
			StructTag(`json:"code"`).
			Annotations(entproto.Field(3)),
		field.String("logo_url").
			StructTag(`json:"logo_url"`).
			Annotations(entproto.Field(4)).
			Optional().
			Nillable(),
		field.String("address").
			StructTag(`json:"address"`).
			Annotations(entproto.Field(5)).
			Optional().
			Nillable(),
		field.String("phone").
			StructTag(`json:"phone"`).
			Annotations(entproto.Field(6)).
			Optional().
			Nillable(),
		field.String("email").
			StructTag(`json:"email"`).
			Annotations(entproto.Field(7)).
			Optional().
			Nillable(),
		field.String("website").
			StructTag(`json:"website"`).
			Annotations(entproto.Field(8)).
			Optional().
			Nillable(),
		field.Time("created_at").
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(9)).
			Default(time.Now),
		field.Time("updated_at").
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(10)).
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Int("parent_id").
			Optional().
			Nillable().
			StructTag(`json:"parent_id"`).
			Annotations(entproto.Field(11)),
	}
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("parent", Organization.Type).
			Ref("children").
			Unique().
			Field("parent_id").
			Annotations(entproto.Field(12)),
		edge.To("children", Organization.Type).
			Annotations(entproto.Field(13)),
		edge.To("departments", Department.Type).
			Annotations(entproto.Field(14)),
		edge.To("projects", Project.Type).
			Annotations(entproto.Field(15)),
		edge.To("labels", Label.Type).
			Annotations(entproto.Field(16)),
		edge.To("leave_requests", LeaveRequest.Type).
			Annotations(entproto.Field(17)),
	}
}

func (Organization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
