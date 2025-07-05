package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DrainedResult holds the schema definition for the DrainedResult entity.
type DrainedResult struct {
	ent.Schema
}

// Fields of the DrainedResult.
func (DrainedResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int("drained_percent"),
		field.Int("avg_passing_score"),
		field.Int("min_passing_score"),
		field.Int("max_passing_score"),
		field.Int("med_passing_score"),
		field.Int("avg_last_admitted_rating_place"),
		field.Int("min_last_admitted_rating_place"),
		field.Int("max_last_admitted_rating_place"),
		field.Int("med_last_admitted_rating_place"),
		field.Int("run_id"),
	}
}

// Edges of the DrainedResult.
func (DrainedResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("heading", Heading.Type).
			Ref("drained_results").
			Unique().
			Required(),
		edge.To("run", Run.Type).
			Unique().
			Required().
			Field("run_id"),
	}
}

// Indexes of the DrainedResult.
func (DrainedResult) Indexes() []ent.Index {
	return []ent.Index{
		// Index for run-based queries
		index.Fields("run_id"),
		// Composite index for run + drained_percent queries (used heavily in results API)
		index.Fields("run_id", "drained_percent"),
		// Index for drained_percent queries (used in steps calculation)
		index.Fields("drained_percent"),
	}
}
