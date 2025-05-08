package book

import (
	"context"

	"github.com/lunn06/library/bookinfo/internal/app/repository/errors"
	"github.com/lunn06/library/bookinfo/internal/domain"
	"github.com/lunn06/library/bookinfo/internal/infrastructure/db/ent"
	"github.com/lunn06/library/bookinfo/internal/infrastructure/db/ent/book"
	"github.com/lunn06/library/bookinfo/internal/infrastructure/db/ent/converter"
)

type Repo interface {
	Get(ctx context.Context, id int) (domain.Book, error)
	SearchByTitleWithLimitOffset(ctx context.Context, title string, limit int, offset int) ([]domain.Book, error)
	Put(ctx context.Context, book domain.Book, authorsIDs []int, genresIDs []int) (domain.Book, error)
	Update(ctx context.Context, book domain.Book, authorsIDs []int, genresIDs []int) error
	Delete(ctx context.Context, id int) error
}

var _ Repo = (*EntRepo)(nil)

func NewEntRepo(client *ent.Client) *EntRepo {
	return &EntRepo{BookClient: client.Book}
}

type EntRepo struct {
	*ent.BookClient
}

// TODO: add errors handling

func (ebr *EntRepo) Get(ctx context.Context, id int) (domain.Book, error) {
	entBook, err := ebr.BookClient.Get(ctx, id)
	if ent.IsNotFound(err) {
		return domain.Book{}, errors.ErrNotFound{Inner: err}
	}
	if err != nil {
		return domain.Book{}, err
	}

	return converter.BookToDomain(entBook), nil
}

func (ebr *EntRepo) SearchByTitleWithLimitOffset(
	ctx context.Context,
	title string,
	limit int,
	offset int,
) ([]domain.Book, error) {
	entBooks, err := ebr.
		Query().
		Where(
			book.Or(
				book.TitleContains(title),
				book.DescriptionContains(title),
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

	return converter.BooksToDomain(entBooks), nil
}

func (ebr *EntRepo) Put(
	ctx context.Context,
	book domain.Book,
	authorsIDs []int,
	genresIDs []int,
) (domain.Book, error) {
	entBook, err := ebr.
		Create().
		SetTitle(book.Title).
		SetDescription(book.Description).
		SetUserID(book.UserID).
		SetBookURL(book.BookURL).
		SetNillableCoverURL(book.CoverURL).
		AddAuthorIDs(authorsIDs...).
		AddGenreIDs(genresIDs...).
		Save(ctx)
	if err != nil {
		return domain.Book{}, err
	}

	return converter.BookToDomain(entBook), nil
}

func (ebr *EntRepo) Update(
	ctx context.Context,
	book domain.Book,
	authorsIDs []int,
	genresIDs []int,
) error {
	err := ebr.
		UpdateOneID(book.ID).
		SetUserID(book.UserID).
		SetTitle(book.Title).
		SetDescription(book.Description).
		SetBookURL(book.BookURL).
		SetNillableCoverURL(book.CoverURL).
		AddAuthorIDs(authorsIDs...).
		AddGenreIDs(genresIDs...).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return errors.ErrNotFound{Inner: err}
	}

	return err
}

func (ebr *EntRepo) Delete(ctx context.Context, id int) error {
	err := ebr.
		DeleteOneID(id).
		Exec(ctx)

	return err
}
