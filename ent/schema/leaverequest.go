package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LeaveRequest holds the schema definition for the LeaveRequest entity.
type LeaveRequest struct {
	ent.Schema
}

// Fields of the LeaveRequest.
func (LeaveRequest) Fields() []ent.Field {
	return []ent.Field{
		field.Float("total_days").
			Annotations(entproto.Field(2)),
		field.Time("start_at").
			Annotations(entproto.Field(3)),
		field.Time("end_at").
			Annotations(entproto.Field(4)),
		field.String("reason").
			Annotations(entproto.Field(5)).
			Optional().
			Nillable(),
		field.Enum("type").
			Values("annual", "unpaid").
			Default("annual").
			Annotations(
				entproto.Field(6),
				entproto.Enum(map[string]int32{
					"annual": 0,
					"unpaid": 1,
				}),
			),
		field.Enum("status").
			Values("pending", "rejected", "approved").
			Default("pending").
			Annotations(
				entproto.Field(7),
				entproto.Enum(map[string]int32{
					"pending":  0,
					"rejected": 1,
					"approved": 2,
				}),
			),
		field.Int("org_id").
			Annotations(entproto.Field(8)),
		field.Int("employee_id").
			Annotations(entproto.Field(9)),
		field.Time("created_at").
			Annotations(entproto.Field(10)).
			Default(time.Now),
		field.Time("updated_at").
			Annotations(entproto.Field(11)).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the LeaveRequest.
func (LeaveRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("leave_approves", LeaveApproval.Type).
			Annotations(entproto.Field(12)),
		edge.From("applicant", Employee.Type).
			Ref("leave_requests").
			Field("employee_id").
			Required().
			Annotations(entproto.Field(13)).
			Unique(),
		edge.From("organization", Organization.Type).
			Ref("leave_requests").
			Field("org_id").
			Required().
			Annotations(entproto.Field(14)).
			Unique(),
	}
}

func (LeaveRequest) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
