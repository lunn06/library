package author

import (
	"context"

	"github.com/lunn06/library/book/internal/app/repository/errors"
	"github.com/lunn06/library/book/internal/domain"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent/author"
	"github.com/lunn06/library/book/internal/infrastructure/db/ent/converter"
)

type Repo interface {
	Get(ctx context.Context, id int) (domain.Author, error)
	SearchByNameWithLimitOffset(ctx context.Context, name string, limit int, offset int) ([]domain.Author, error)
	Put(ctx context.Context, author domain.Author, booksIDs ...int) (domain.Author, error)
	Update(ctx context.Context, author domain.Author, booksIDs ...int) error
	Delete(ctx context.Context, id int) error
}

var _ Repo = (*EntRepo)(nil)

func NewEntRepo(client *ent.Client) *EntRepo {
	return &EntRepo{AuthorClient: client.Author}
}

type EntRepo struct {
	*ent.AuthorClient
}

func (ear *EntRepo) Get(ctx context.Context, id int) (domain.Author, error) {
	entAuthor, err := ear.AuthorClient.Get(ctx, id)
	if ent.IsNotFound(err) {
		return domain.Author{}, errors.ErrNotFound{Inner: err}
	}
	if err != nil {
		return domain.Author{}, err
	}

	return converter.AuthorToDomain(entAuthor), nil
}

func (ear *EntRepo) SearchByNameWithLimitOffset(
	ctx context.Context,
	name string,
	limit int,
	offset int,
) ([]domain.Author, error) {
	entAuthors, err := ear.
		Query().
		Where(
			author.Or(
				author.NameContains(name),
				author.DescriptionContains(name),
			),
		).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if ent.IsNotFound(err) {
		return nil, errors.ErrNotFound{Inner: err}
	}
	if err != nil {
		return nil, err
	}

	return converter.AuthorsToDomain(entAuthors), nil
}

func (ear *EntRepo) Put(ctx context.Context, author domain.Author, booksIDs ...int) (domain.Author, error) {
	entAuthor, err := ear.
		Create().
		SetName(author.Name).
		SetDescription(author.Description).
		AddBookIDs(booksIDs...).
		Save(ctx)
	if err != nil {
		return domain.Author{}, err
	}

	return converter.AuthorToDomain(entAuthor), nil
}

func (ear *EntRepo) Update(ctx context.Context, author domain.Author, booksIDs ...int) error {
	err := ear.AuthorClient.
		UpdateOneID(author.ID).
		SetName(author.Name).
		SetDescription(author.Description).
		AddBookIDs(booksIDs...).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return errors.ErrNotFound{Inner: err}
	}

	return err
}

func (ear *EntRepo) Delete(ctx context.Context, id int) error {
	err := ear.AuthorClient.
		DeleteOneID(id).
		Exec(ctx)

	return err
}
