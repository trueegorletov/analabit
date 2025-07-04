package schema

import (
	"analabit/core"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Application holds the schema definition for the Application entity.
type Application struct {
	ent.Schema
}

// Fields of the Application.
func (Application) Fields() []ent.Field {
	return []ent.Field{
		field.String("student_id"),
		field.Int("priority"),
		field.Int("competition_type").GoType(core.Competition(0)),
		field.Int("rating_place"),
		field.Int("score"),
		field.Int("run_id"),
		field.Bool("original_submitted").Default(false),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Application.
func (Application) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("heading", Heading.Type).
			Ref("applications").
			Unique().
			Required(),
		edge.To("run", Run.Type).
			Unique().
			Required().
			Field("run_id"),
	}
}

// Indexes of the Application.
func (Application) Indexes() []ent.Index {
	return []ent.Index{
		// Index for run-based queries
		index.Fields("run_id"),
		// Composite index for run + student queries (used in API handlers)
		index.Fields("run_id", "student_id"),
		// Index for original_submitted queries
		index.Fields("original_submitted"),
	}
}
