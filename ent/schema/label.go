package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Label holds the schema definition for the Label entity.
type Label struct {
	ent.Schema
}

// Fields of the Label.
func (Label) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(entproto.Field(2)).
			NotEmpty(),
		field.String("description").
			Annotations(entproto.Field(3)).
			Optional(),
		field.String("color").
			Annotations(entproto.Field(4)).
			NotEmpty(),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(5)),
		field.Time("updated_at").
			Annotations(entproto.Field(6)).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Label.
func (Label) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Annotations(entproto.Field(7)),
	}
}

func (Label) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
