package converter

import (
	"github.com/lunn06/library/book/internal/domain"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent"
)

func AuthorsToDomain(entAuthors ent.Authors) []domain.Author {
	authors := make([]domain.Author, len(entAuthors))
	for i, entAuthor := range entAuthors {
		authors[i] = AuthorToDomain(entAuthor)
	}

	return authors
}

func AuthorsToEnt(authors []domain.Author) ent.Authors {
	entAuthors := make([]*ent.Author, len(authors))
	for i, author := range authors {
		entAuthors[i] = AuthorToEnt(author)
	}

	return entAuthors
}

func AuthorToDomain(author *ent.Author) domain.Author {
	return domain.Author{
		ID:          author.ID,
		Name:        author.Name,
		Description: author.Description,
		Books:       BooksToDomain(author.Edges.Books),
	}
}

func AuthorToEnt(author domain.Author) *ent.Author {
	return &ent.Author{
		ID:          author.ID,
		Name:        author.Name,
		Description: author.Description,
		Edges: ent.AuthorEdges{
			Books: BooksToEnt(author.Books),
		},
	}
}
