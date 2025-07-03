package schema

import (
	"analabit/core"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.Int("iteration"),
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
	}
}
