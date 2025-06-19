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

// ZaloEmployee holds the schema definition for the ZaloEmployee entity.
type ZaloEmployee struct {
	ent.Schema
}

// Fields of the ZaloEmployee.
func (ZaloEmployee) Fields() []ent.Field {
	return []ent.Field{
		field.String("zalo_uid").
			NotEmpty().
			StructTag(`json:"zalo_uid"`).
			Annotations(entproto.Field(2)),
		field.Int("employee_id").
			StructTag(`json:"employee_id"`).
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

// Edges of the ZaloEmployee.
func (ZaloEmployee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("employee", Employee.Type).
			Ref("zalo_employee").
			Field("employee_id").
			Unique().
			Required().
			StructTag(`json:"employee"`).
			Annotations(entproto.Field(6)),
	}
}

// Indexes of the ZaloEmployee.
func (ZaloEmployee) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint: mỗi employee_id chỉ tương ứng với 1 zalo_uid
		index.Fields("employee_id").
			Unique(),
		// Index cho zalo_uid để tìm kiếm nhanh
		index.Fields("zalo_uid"),
	}
}

// Annotations of the ZaloEmployee.
func (ZaloEmployee) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
	}
}
