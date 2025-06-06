package schema

import (
	"analabit/core"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
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
