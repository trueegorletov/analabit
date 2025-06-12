package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DrainedResult holds the schema definition for the DrainedResult entity.
type DrainedResult struct {
	ent.Schema
}

// Fields of the DrainedResult.
func (DrainedResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int("drained_percent"),
		field.Int("passing_score"),
		field.Int("last_admitted_rating_place"),
		field.Int("iteration"),
	}
}

// Edges of the DrainedResult.
func (DrainedResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("heading", Heading.Type).
			Ref("drained_results").
			Unique().
			Required(),
	}
}
