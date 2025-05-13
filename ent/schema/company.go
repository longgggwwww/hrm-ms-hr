package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Company struct {
	ent.Schema
}

func (Company) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("name").
			Annotations(entproto.Field(2)),
		field.String("address").
			Annotations(entproto.Field(3)),
		field.Time("created_at").
			Default(time.Now).
			Annotations(entproto.Field(4)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(5)),
	}
}

func (Company) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("branches", Branch.Type),
	}
}
