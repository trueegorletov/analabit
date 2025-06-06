package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
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
	}
}
