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

type Branch struct {
	ent.Schema
}

func (Branch) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("name").
			Annotations(entproto.Field(2)),
		field.String("code").
			Unique().
			Annotations(entproto.Field(3)),
		field.String("address").
			Annotations(entproto.Field(4)),
		field.String("contact_info").
			Annotations(entproto.Field(5)),
		field.Time("created_at").
			Default(time.Now).
			Annotations(entproto.Field(6)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(7)),
	}
}

func (Branch) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("company", Company.Type).
			Ref("branches").
			Unique().
			Annotations(entproto.Field(8)),
	}
}

func (Branch) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
