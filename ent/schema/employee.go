package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Employee holds the schema definition for the Employee entity.
type Employee struct {
	ent.Schema
}

// Fields of the Employee.
func (Employee) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			Unique(). // mã nhân viên yêu cầu uniq kể cả khác tổ chức
			Optional().
			Annotations(entproto.Field(2)),
		field.String("code").
			NotEmpty().
			Unique().
			Annotations(entproto.Field(3)),
		field.Enum("status").
			Values("active", "inactive").
			Default("active").
			Annotations(entproto.Field(4),
				entproto.Enum(map[string]int32{
					"active":   0,
					"inactive": 1,
				}),
			),
		field.Int("position_id").
			Annotations(entproto.Field(5)),
		field.Time("joining_at").
			Annotations(entproto.Field(6)),
		field.Int("org_id").
			Annotations(entproto.Field(7)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(8)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(9)),
	}
}

// Edges of the Employee.
func (Employee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("position", Position.Type).
			Ref("employees").
			Field("position_id").
			Unique().
			Required().
			Annotations(entproto.Field(10)),
	}
}

func (Employee) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
