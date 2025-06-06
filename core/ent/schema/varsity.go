package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Varsity holds the schema definition for the Varsity entity.
type Varsity struct {
	ent.Schema
}

// Fields of the Varsity.
func (Varsity) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").Unique(),
		field.String("name"),
	}
}

// Edges of the Varsity.
func (Varsity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("headings", Heading.Type),
	}
}
