package converter

import (
	"github.com/lunn06/library/book/internal/domain"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent"
)

func GenresToDomain(entGenres ent.Genres) []domain.Genre {
	genres := make([]domain.Genre, len(entGenres))
	for i, entGenre := range entGenres {
		genres[i] = GenreToDomain(entGenre)
	}

	return genres
}

func GenresToEnt(genres []domain.Genre) ent.Genres {
	entGenres := make([]*ent.Genre, len(genres))
	for i, genre := range genres {
		entGenres[i] = GenreToEnt(genre)
	}

	return entGenres
}

func GenreToDomain(genre *ent.Genre) domain.Genre {
	return domain.Genre{
		ID:          genre.ID,
		Title:       genre.Title,
		Description: genre.Description,
		Books:       BooksToDomain(genre.Edges.Books),
	}
}

func GenreToEnt(genre domain.Genre) *ent.Genre {
	return &ent.Genre{
		ID:          genre.ID,
		Title:       genre.Title,
		Description: genre.Description,
		Edges: ent.GenreEdges{
			Books: BooksToEnt(genre.Books),
		},
	}
}
