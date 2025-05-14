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

// Position holds the schema definition for the Position entity.
type Position struct {
	ent.Schema
}

// Fields of the Position.
func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("code").
			Unique().
			Annotations(entproto.Field(3)),
		field.UUID("department_id", uuid.UUID{}).
			Annotations(entproto.Field(4)),
		field.UUID("parent_id", uuid.UUID{}).
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

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("employees", Employee.Type).
			Annotations(entproto.Field(8)),
		edge.From("department", Department.Type).
			Ref("positions").
			Field("department_id").
			Unique().
			Required().
			Annotations(entproto.Field(9)),
		// edge.To("children", Position.Type).
		// 	From("parent").
		// 	Unique().
		// 	Field("parent_id").
		// 	Annotations(entproto.Field(8)),
	}
}

func (Position) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
