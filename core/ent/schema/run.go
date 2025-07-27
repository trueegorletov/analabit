package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Run holds the schema definition for the Run entity.
type Run struct {
	ent.Schema
}

// Fields of the Run.
func (Run) Fields() []ent.Field {
	return []ent.Field{
		field.Time("triggered_at").Default(time.Now),
		field.JSON("payload_meta", map[string]any{}).Optional(),
		field.Bool("finished").Default(false),
		field.Time("finished_at").Default(time.Now),
	}
}

// Edges of the Run.
func (Run) Edges() []ent.Edge {
	return []ent.Edge{
		// No edges needed here; child tables point to Run via run_id foreign key
	}
}
