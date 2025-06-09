package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type AppointmentHistory struct {
	ent.Schema
}

// Fields of the AppointmentHistory.
func (AppointmentHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("employee_id").
			StructTag(`json:"employee_id"`).
			Annotations(entproto.Field(2)),
		field.String("position_name").
			StructTag(`json:"position_name"`).
			Annotations(entproto.Field(3)),
		field.Time("joining_at").
			StructTag(`json:"joining_at"`).
			Annotations(entproto.Field(4)),
		field.String("description").Optional().
			StructTag(`json:"description"`).
			Annotations(entproto.Field(5)),
		field.Strings("attachment_urls").Optional().
			StructTag(`json:"attachment_urls"`).
			Annotations(entproto.Field(6)),
		field.Time("created_at").Default(time.Now).
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(7)),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(8)),
	}
}

func (AppointmentHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("employee", Employee.Type).
			Ref("appointment_histories").
			Field("employee_id").
			Unique().
			Required().
			Annotations(entproto.Field(9), entsql.OnDelete(entsql.Cascade)),
	}
}

func (AppointmentHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
