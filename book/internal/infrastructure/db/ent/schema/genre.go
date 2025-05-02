package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Genre holds the schema definition for the Genre entity.
type Genre struct {
	ent.Schema
}

// Fields of the Genre.
func (Genre) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StructTag(`json:"id,omitempty"`),
		field.String("title").
			Unique().
			NotEmpty(),
		field.String("description").
			NotEmpty(),
	}
}

// Edges of the Genre.
func (Genre) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("books", Book.Type).
			Ref("genres"),
	}
}
