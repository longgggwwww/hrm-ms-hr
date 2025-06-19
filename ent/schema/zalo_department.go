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

// ZaloDepartment holds the schema definition for the ZaloDepartment entity.
type ZaloDepartment struct {
	ent.Schema
}

// Fields of the ZaloDepartment.
func (ZaloDepartment) Fields() []ent.Field {
	return []ent.Field{
		field.String("group_id").
			NotEmpty().
			StructTag(`json:"group_id"`).
			Annotations(entproto.Field(2)),
		field.Int("department_id").
			StructTag(`json:"department_id"`).
			Annotations(entproto.Field(3)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(4)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(5)),
	}
}

// Edges of the ZaloDepartment.
func (ZaloDepartment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("department", Department.Type).
			Ref("zalo_department").
			Field("department_id").
			Unique().
			Required().
			StructTag(`json:"department"`).
			Annotations(entproto.Field(6)),
	}
}

// Indexes of the ZaloDepartment.
func (ZaloDepartment) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint: mỗi department_id chỉ tương ứng với 1 group_id
		index.Fields("department_id").
			Unique(),
		// Index cho group_id để tìm kiếm nhanh
		index.Fields("group_id"),
	}
}

// Annotations of the ZaloDepartment.
func (ZaloDepartment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
	}
}
