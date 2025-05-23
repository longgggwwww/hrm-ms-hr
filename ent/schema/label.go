package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Label holds the schema definition for the Label entity.
type Label struct {
	ent.Schema
}

// Fields of the Label.
func (Label) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("code").
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(3)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(4)),
		field.String("color").
			NotEmpty().
			Annotations(entproto.Field(5)),
		field.UUID("project_id", uuid.UUID{}).
			Annotations(entproto.Field(6)),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(7)),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(8)),
	}
}

// Edges of the Label.
func (Label) Edges() []ent.Edge {
	return []ent.Edge{
		// n-1 with Project
		edge.From("project", Project.Type).
			Ref("labels").
			Field("project_id").
			Unique().
			Required().
			Annotations(entproto.Field(9)),
		// n-n with Task
		edge.To("tasks", Task.Type).
			Annotations(entproto.Field(10)),
	}
}

func (Label) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code", "project_id").Unique(),
	}
}

func (Label) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
