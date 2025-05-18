package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
		field.Time("created_at").
			Annotations(entproto.Field(3)).
			Default(time.Now),
		field.Time("updated_at").
			Annotations(entproto.Field(4)).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the LeaveApproval.
func (LeaveApproval) Edges() []ent.Edge {
	return nil
}

func (LeaveApproval) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}