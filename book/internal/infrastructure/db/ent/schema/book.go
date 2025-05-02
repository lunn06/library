package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StructTag(`json:"id,omitempty"`),
		field.Int("user_id"),
		field.String("title").
			NotEmpty(),
		field.String("description").
			NotEmpty(),
		field.String("book_url").
			NotEmpty(),
		field.String("cover_url").
			Optional().
			Nillable(),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("genres", Genre.Type),
		edge.To("authors", Author.Type),
	}
}
