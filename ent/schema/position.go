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
			StructTag(`json:"name"`).
			Annotations(entproto.Field(2)),
		field.String("code").
			StructTag(`json:"code"`).
			Annotations(entproto.Field(3)),
		field.Int("department_id").
			StructTag(`json:"department_id"`).
			Annotations(entproto.Field(4)),
		field.Int("parent_id").
			Optional().
			StructTag(`json:"parent_id"`).
			Annotations(entproto.Field(5)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(6)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(7)),
	}
}

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("employees", Employee.Type).
			StructTag(`json:"employees"`).
			Annotations(entproto.Field(8)),
		edge.From("departments", Department.Type). // Lỗi chính tả
								Ref("positions").
								Field("department_id").
								Unique().
								Required().
								StructTag(`json:"department"`).
								Annotations(entproto.Field(9)),

		// edge đến con
		edge.To("children", Position.Type).
			StructTag(`json:"children"`).
			Annotations(entproto.Field(10)), // QUAN TRỌNG: để entproto generate RPC field

		// edge đến cha
		edge.From("parent", Position.Type).
			Ref("children").
			Unique().
			Field("parent_id").
			StructTag(`json:"parent"`).
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
