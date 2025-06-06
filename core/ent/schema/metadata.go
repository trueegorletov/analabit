package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Metadata holds the schema definition for the Metadata entity.
type Metadata struct {
	ent.Schema
}

// Fields of the Metadata.
func (Metadata) Fields() []ent.Field {
	return []ent.Field{
		field.Int("last_applications_iteration"),
		field.Int("last_calculations_iteration"),
		field.Bool("uploading_lock"),
	}
}

// Edges of the Metadata.
func (Metadata) Edges() []ent.Edge {
	return nil
}
