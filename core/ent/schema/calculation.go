package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Calculation holds the schema definition for the Calculation entity.
type Calculation struct {
	ent.Schema
}

// Fields of the Calculation.
func (Calculation) Fields() []ent.Field {
	return []ent.Field{
		field.String("student_id"),
		field.Int("admitted_place"),
		field.Int("iteration"),
		field.Int("run_id"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Calculation.
func (Calculation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("heading", Heading.Type).
			Ref("calculations").
			Unique().
			Required(),
		edge.To("run", Run.Type).
			Unique().
			Required().
			Field("run_id"),
	}
}

// Indexes of the Calculation.
func (Calculation) Indexes() []ent.Index {
	return []ent.Index{
		// Index for run-based queries
		index.Fields("run_id"),
		// Composite index for run + student queries (used in API handlers)
		index.Fields("run_id", "student_id"),
		// Index for iteration-based queries (backward compatibility)
		index.Fields("iteration"),
		// Composite index for student + iteration (legacy queries)
		index.Fields("student_id", "iteration"),
		// Index for admitted_place queries (used in results calculation)
		index.Fields("admitted_place"),
	}
}
