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

type Company struct {
	ent.Schema
}

func (Company) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("code").
			NotEmpty().
			Unique().
			Annotations(entproto.Field(3)),
		field.String("address").
			Optional().
			Annotations(entproto.Field(5)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(6)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(7)),
	}
}

func (Company) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("branches", Branch.Type).
			Annotations(entproto.Field(8)),
	}
}

func (Company) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
