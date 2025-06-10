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
			StructTag(`json:"total_days"`).
			Annotations(entproto.Field(2)),
		field.Time("start_at").
			StructTag(`json:"start_at"`).
			Annotations(entproto.Field(3)),
		field.Time("end_at").
			StructTag(`json:"end_at"`).
			Annotations(entproto.Field(4)),
		field.String("reason").
			StructTag(`json:"reason"`).
			Annotations(entproto.Field(5)).
			Optional().
			Nillable(),
		field.Enum("type").
			Values("annual", "unpaid").
			Default("annual").
			StructTag(`json:"type"`).
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
			StructTag(`json:"status"`).
			Annotations(
				entproto.Field(7),
				entproto.Enum(map[string]int32{
					"pending":  0,
					"rejected": 1,
					"approved": 2,
				}),
			),
		field.Int("org_id").
			StructTag(`json:"org_id"`).
			Annotations(entproto.Field(8)),
		field.Int("employee_id").
			StructTag(`json:"employee_id"`).
			Annotations(entproto.Field(9)),
		field.Time("created_at").
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(10)).
			Default(time.Now),
		field.Time("updated_at").
			StructTag(`json:"updated_at"`).
			Annotations(entproto.Field(11)).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the LeaveRequest.
func (LeaveRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("leave_approves", LeaveApproval.Type).
			StructTag(`json:"leave_approves"`).
			Annotations(entproto.Field(12)),
		edge.From("applicant", Employee.Type).
			Ref("leave_requests").
			Field("employee_id").
			Required().
			StructTag(`json:"applicant"`).
			Annotations(entproto.Field(13)).
			Unique(),
		edge.From("organization", Organization.Type).
			Ref("leave_requests").
			Field("org_id").
			Required().
			StructTag(`json:"organization"`).
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
