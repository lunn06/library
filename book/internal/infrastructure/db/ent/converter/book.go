package converter

import (
	"github.com/lunn06/library/book/internal/domain"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent"
)

func BooksToDomain(entBooks ent.Books) []domain.Book {
	books := make([]domain.Book, len(entBooks))
	for i, entBook := range entBooks {
		books[i] = BookToDomain(entBook)
	}

	return books
}

func BooksToEnt(books []domain.Book) ent.Books {
	entBooks := make([]*ent.Book, len(books))
	for i, book := range books {
		entBooks[i] = BookToEnt(book)
	}

	return entBooks
}

func BookToDomain(book *ent.Book) domain.Book {
	return domain.Book{
		ID:          book.ID,
		UserID:      book.UserID,
		Title:       book.Title,
		Description: book.Description,
		BookURL:     book.BookURL,
		CoverURL:    book.CoverURL,
		Authors:     AuthorsToDomain(book.Edges.Authors),
		Genres:      GenresToDomain(book.Edges.Genres),
	}
}

func BookToEnt(book domain.Book) *ent.Book {
	return &ent.Book{
		ID:          book.ID,
		UserID:      book.UserID,
		Title:       book.Title,
		Description: book.Description,
		BookURL:     book.BookURL,
		CoverURL:    book.CoverURL,
		Edges: ent.BookEdges{
			Genres:  GenresToEnt(book.Genres),
			Authors: AuthorsToEnt(book.Authors),
		},
	}
}
