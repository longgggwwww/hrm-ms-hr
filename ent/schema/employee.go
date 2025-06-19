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
			StructTag(`json:"user_id"`).
			Annotations(entproto.Field(2)),
		field.String("code").
			NotEmpty().
			Unique().
			StructTag(`json:"code"`).
			Annotations(entproto.Field(3)),
		field.Enum("status").
			Values("active", "inactive").
			Default("active").
			StructTag(`json:"status"`).
			Annotations(entproto.Field(4),
				entproto.Enum(map[string]int32{
					"active":   0,
					"inactive": 1,
				}),
			),
		field.Int("position_id").
			StructTag(`json:"position_id"`).
			Annotations(entproto.Field(5)),
		field.Time("joining_at").
			StructTag(`json:"joining_at"`).
			Annotations(entproto.Field(6)),
		field.Int("org_id").
			StructTag(`json:"org_id"`).
			Annotations(entproto.Field(7)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			StructTag(`json:"created_at"`).
			Annotations(entproto.Field(8)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updated_at"`).
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
			StructTag(`json:"position"`).
			Annotations(entproto.Field(10)),
		edge.To("created_projects", Project.Type).
			StructTag(`json:"created_projects"`).
			Annotations(entproto.Field(11)),
		edge.To("updated_projects", Project.Type).
			StructTag(`json:"updated_projects"`).
			Annotations(entproto.Field(12)),
		edge.From("assigned_tasks", Task.Type).
			Ref("assignees").
			StructTag(`json:"assigned_tasks"`).
			Annotations(entproto.Field(13)),
		edge.To("leave_approves", LeaveApproval.Type).
			StructTag(`json:"leave_approves"`).
			Annotations(entproto.Field(14)),
		edge.To("leave_requests", LeaveRequest.Type).
			StructTag(`json:"leave_requests"`).
			Annotations(entproto.Field(15)),
		edge.To("task_reports", TaskReport.Type).
			StructTag(`json:"task_reports"`).
			Annotations(entproto.Field(16)), // Edge đến TaskReport
		edge.From("projects", Project.Type).
			Ref("members").
			StructTag(`json:"projects"`).
			Annotations(entproto.Field(17)), // Inverse edge for project members
		edge.To("appointment_histories", AppointmentHistory.Type).
			StructTag(`json:"appointment_histories"`).
			Annotations(entproto.Field(18)),
		edge.To("zalo_employee", ZaloEmployee.Type).
			Unique().
			StructTag(`json:"zalo_employee"`).
			Annotations(entproto.Field(19)),
	}
}

func (Employee) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
