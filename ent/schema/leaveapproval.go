package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LeaveApproval holds the schema definition for the LeaveApproval entity.
type LeaveApproval struct {
	ent.Schema
}

// Fields of the LeaveApproval.
func (LeaveApproval) Fields() []ent.Field {
	return []ent.Field{
		field.String("comment").
			Annotations(entproto.Field(2)).
			Optional().
			Nillable(),
		field.Int("leave_request_id").
			Annotations(entproto.Field(3)),
		field.Int("reviewer_id").
			Annotations(entproto.Field(4)),
		field.Time("created_at").
			Annotations(entproto.Field(5)).
			Default(time.Now),
		field.Time("updated_at").
			Annotations(entproto.Field(6)).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the LeaveApproval.
func (LeaveApproval) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("leave_request", LeaveRequest.Type).
			Ref("leave_approves").
			Field("leave_request_id").
			Required().
			Annotations(entproto.Field(7)).
			Unique(),
		edge.From("reviewer", Employee.Type).
			Ref("leave_approves").
			Field("reviewer_id").
			Required().
			Annotations(entproto.Field(8)).
			Unique(),
	}
}

func (LeaveApproval) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
