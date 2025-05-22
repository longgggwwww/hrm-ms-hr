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

// Position holds the schema definition for the Position entity.
type Position struct {
	ent.Schema
}

// Fields of the Position.
func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("code").
			Annotations(entproto.Field(3)),
		field.Int("department_id").
			Annotations(entproto.Field(4)),
		field.Int("parent_id").
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

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("employees", Employee.Type).
			Annotations(entproto.Field(8)),
		edge.From("departments", Department.Type).
			Ref("positions").
			Field("department_id").
			Unique().
			Required().
			Annotations(entproto.Field(9)),

		// edge đến con
		edge.To("children", Position.Type).
			Annotations(entproto.Field(10)), // QUAN TRỌNG: để entproto generate RPC field

		// edge đến cha
		edge.From("parent", Position.Type).
			Ref("children").
			Unique().
			Field("parent_id").
			Annotations(entproto.Field(11)),
	}
}

func (Position) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}

func (Position) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("department_id", "code").Unique(),
	}
}
