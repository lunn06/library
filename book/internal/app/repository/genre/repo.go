package genre

import (
	"context"

	"github.com/lunn06/library/book/internal/domain"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent/converter"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent/genre"
)

type Repo interface {
	Get(ctx context.Context, id int) (domain.Genre, error)
	SearchByTitleWithLimitOffset(ctx context.Context, title string, limit int, offset int) ([]domain.Genre, error)
	Put(ctx context.Context, genre domain.Genre, booksIDs ...int) (domain.Genre, error)
	Update(ctx context.Context, genre domain.Genre, bookIDs ...int) error
	Delete(ctx context.Context, id int) error
}

var _ Repo = (*EntRepo)(nil)

func NewEntRepo(client *ent.Client) *EntRepo {
	return &EntRepo{GenreClient: client.Genre}
}

type EntRepo struct {
	*ent.GenreClient
}

func (egr *EntRepo) Get(ctx context.Context, id int) (domain.Genre, error) {
	entGenre, err := egr.GenreClient.Get(ctx, id)
	if err != nil {
		return domain.Genre{}, err
	}

	return converter.GenreToDomain(entGenre), nil
}

func (egr *EntRepo) SearchByTitleWithLimitOffset(
	ctx context.Context,
	title string,
	limit int,
	offset int,
) ([]domain.Genre, error) {
	entGenres, err := egr.
		Query().
		Where(
			genre.Or(
				genre.Title(title),
				genre.DescriptionContains(title),
			),
		).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return converter.GenresToDomain(entGenres), nil
}

func (egr *EntRepo) Put(ctx context.Context, genre domain.Genre, booksIDs ...int) (domain.Genre, error) {
	entGenre, err := egr.
		Create().
		SetTitle(genre.Title).
		SetDescription(genre.Description).
		AddBookIDs(booksIDs...).
		Save(ctx)
	if err != nil {
		return domain.Genre{}, err
	}

	return converter.GenreToDomain(entGenre), nil
}

func (egr *EntRepo) Update(ctx context.Context, genre domain.Genre, booksIDs ...int) error {
	err := egr.
		UpdateOneID(genre.ID).
		SetTitle(genre.Title).
		SetDescription(genre.Description).
		AddBookIDs(booksIDs...).
		Exec(ctx)

	return err
}

func (egr *EntRepo) Delete(ctx context.Context, id int) error {
	err := egr.GenreClient.
		DeleteOneID(id).
		Exec(ctx)

	return err
}
