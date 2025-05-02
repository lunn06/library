package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Author holds the schema definition for the Author entity.
type Author struct {
	ent.Schema
}

// Fields of the Author.
func (Author) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StructTag(`json:"id,omitempty"`),
		field.String("name").
			NotEmpty(),
		field.String("description").
			NotEmpty(),
	}
}

// Edges of the Author.
func (Author) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("books", Book.Type).
			Ref("authors"),
	}
}
