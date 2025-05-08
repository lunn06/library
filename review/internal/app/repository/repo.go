package repository

import (
	"context"

	"github.com/lunn06/library/review/internal/app/repository/errors"
	"github.com/lunn06/library/review/internal/domain"
	"github.com/lunn06/library/review/internal/infrastructure/db/ent"
	"github.com/lunn06/library/review/internal/infrastructure/db/ent/converter"
	"github.com/lunn06/library/review/internal/infrastructure/db/ent/review"
)

type ReviewRepo interface {
	GetAllByBookID(ctx context.Context, bookID int) ([]domain.Review, error)
	Get(ctx context.Context, id int) (domain.Review, error)
	Put(ctx context.Context, book domain.Review) (domain.Review, error)
	Update(ctx context.Context, book domain.Review) error
	Delete(ctx context.Context, id int) error
}

var _ ReviewRepo = (*EntReviewRepo)(nil)

func NewEntReviewRepo(client *ent.Client) *EntReviewRepo {
	return &EntReviewRepo{ReviewClient: client.Review}
}

type EntReviewRepo struct {
	*ent.ReviewClient
}

// TODO: add errors handling

func (rr *EntReviewRepo) GetAllByBookID(ctx context.Context, bookID int) ([]domain.Review, error) {
	entReviews, err := rr.ReviewClient.
		Query().
		Where(
			review.BookID(bookID),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return converter.ReviewsToDomain(entReviews), nil
}

func (rr *EntReviewRepo) Get(ctx context.Context, id int) (domain.Review, error) {
	entBook, err := rr.ReviewClient.Get(ctx, id)
	if ent.IsNotFound(err) {
		return domain.Review{}, errors.ErrNotFound{Inner: err}
	}
	if err != nil {
		return domain.Review{}, err
	}

	return converter.ReviewToDomain(entBook), nil
}

func (rr *EntReviewRepo) Put(
	ctx context.Context,
	review domain.Review,
) (domain.Review, error) {
	cte := rr.
		Create().
		SetUserID(review.UserID).
		SetBookID(review.BookID).
		SetTitle(review.Title).
		SetText(review.Text).
		SetScore(review.Score)

	if !review.CreatedAt.IsZero() {
		cte = cte.SetCreatedAt(review.CreatedAt)
	}

	entReview, err := cte.Save(ctx)
	if err != nil {
		return domain.Review{}, err
	}

	return converter.ReviewToDomain(entReview), nil
}

func (rr *EntReviewRepo) Update(
	ctx context.Context,
	review domain.Review,
) error {
	err := rr.
		UpdateOneID(review.ID).
		SetTitle(review.Title).
		SetText(review.Text).
		SetScore(review.Score).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return errors.ErrNotFound{Inner: err}
	}

	return err
}

func (rr *EntReviewRepo) Delete(ctx context.Context, id int) error {
	err := rr.
		DeleteOneID(id).
		Exec(ctx)

	return err
}
