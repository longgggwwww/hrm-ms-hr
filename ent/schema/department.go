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

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("name").
			Annotations(entproto.Field(2)),
		field.String("code").
			Unique().
			Annotations(entproto.Field(3)),
		field.UUID("branch_id", uuid.UUID{}).
			Annotations(entproto.Field(4)),
		field.Time("created_at").
			Default(time.Now).
			Annotations(entproto.Field(5)),
		field.Time("updated_at").
			Default(time.Now).
			Annotations(entproto.Field(6)),
	}
}

// Edges of the Department.
func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("positions", Position.Type).
			Annotations(entproto.Field(7)),
	}
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
