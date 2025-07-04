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

// Label holds the schema definition for the Label entity.
type Label struct {
	ent.Schema
}

// Fields of the Label.
func (Label) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			StructTag(`json:"name"`).
			Annotations(entproto.Field(2)).
			NotEmpty(),
		field.String("description").
			StructTag(`json:"description"`).
			Annotations(entproto.Field(3)).
			Optional(),
		field.String("color").
			StructTag(`json:"color"`).
			Annotations(entproto.Field(4)).
			NotEmpty(),
		field.Int("org_id").
			StructTag(`json:"org_id"`).
			Annotations(entproto.Field(5)).
			Optional(),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(6)),
		field.Time("updated_at").
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(7)).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Label.
func (Label) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tasks", Task.Type).
			Ref("labels").
			StructTag(`json:"tasks"`).
			Annotations(entproto.Field(8)),
		edge.From("organization", Organization.Type).
			Ref("labels").
			Unique().
			Field("org_id").
			StructTag(`json:"organization"`).
			Annotations(entproto.Field(9)),
	}
}

// Indexes of the Label.
func (Label) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_id"),
	}
}

func (Label) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
