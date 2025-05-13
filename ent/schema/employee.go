package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Employee holds the schema definition for the Employee entity.
type Employee struct {
	ent.Schema
}

// Fields of the Employee.
func (Employee) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Annotations(entproto.Field(1)),
		field.String("employee_id").Unique().Annotations(entproto.Field(2)),
		field.String("code").Unique().Annotations(entproto.Field(3)),
		field.Bool("status").Annotations(entproto.Field(4)),
		field.UUID("position_id", uuid.UUID{}).Annotations(entproto.Field(5)),
		field.Time("joining_at").Annotations(entproto.Field(6)),
		field.UUID("branch_id", uuid.UUID{}).Annotations(entproto.Field(7)),
		field.Time("created_at").
			Default(time.Now).Annotations(entproto.Field(8)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).Annotations(entproto.Field(9)),
		field.UUID("department_id", uuid.UUID{}).Annotations(entproto.Field(10)),
	}
}

// Edges of the Employee.
func (Employee) Edges() []ent.Edge {
	return nil
}

func (Employee) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
