package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Position holds the schema definition for the Position entity.
type Position struct {
	ent.Schema
}

// Fields of the Position.
func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Annotations(entproto.Field(1)),
		field.String("name").NotEmpty().Annotations(entproto.Field(2)),
		field.String("code").Unique().NotEmpty().Annotations(entproto.Field(3)),
		field.UUID("department_id", uuid.UUID{}).Annotations(entproto.Field(4)),
		field.UUID("parent_id", uuid.UUID{}).Annotations(entproto.Field(5)),
		field.Time("created_at").Default(time.Now).Immutable().Annotations(entproto.Field(6)),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Annotations(entproto.Field(7)),
	}
}

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return nil
}

func (Position) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
