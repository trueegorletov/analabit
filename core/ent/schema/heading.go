package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Heading holds the schema definition for the Heading entity.
type Heading struct {
	ent.Schema
}

// Fields of the Heading.
func (Heading) Fields() []ent.Field {
	return []ent.Field{
		field.Int("regular_capacity"),
		field.Int("target_quota_capacity"),
		field.Int("dedicated_quota_capacity"),
		field.Int("special_quota_capacity"),
		field.String("code").Unique(),
		field.String("name"),
	}
}

// Edges of the Heading.
func (Heading) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("varsity", Varsity.Type).Ref("headings").Unique().Required(),
		edge.To("applications", Application.Type),
		edge.To("calculations", Calculation.Type),
	}
}
