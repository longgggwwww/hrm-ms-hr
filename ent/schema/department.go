package schema

import "entgo.io/ent"

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return nil
}

// Edges of the Department.
func (Department) Edges() []ent.Edge {
	return nil
}
