package service

import (
	"context"

	"github.com/lunn06/library/review/internal/app/repository"
	repoerrors "github.com/lunn06/library/review/internal/app/repository/errors"
	"github.com/lunn06/library/review/internal/app/service/errors"
	"github.com/lunn06/library/review/internal/domain"
)

// TODO: add errors handling

type UpdateRequest struct {
	ID    int
	Title string
	Text  string
	Score int
}

type CreateRequest struct {
	UserID int
	BookID int
	Title  string
	Text   string
	Score  int
}

func NewReviewService(repo repository.ReviewRepo) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

type ReviewService struct {
	repo repository.ReviewRepo
}

func (s *ReviewService) GetAllByBookID(ctx context.Context, bookID int) ([]domain.Review, error) {
	return s.repo.GetAllByBookID(ctx, bookID)
}

func (s *ReviewService) Create(ctx context.Context, req CreateRequest) (int, error) {
	score, err := domain.NewScore(req.Score)
	if err != nil {
		return 0, err
	}
	review := domain.Review{
		UserID: req.UserID,
		BookID: req.BookID,
		Title:  req.Title,
		Text:   req.Text,
		Score:  score,
	}

	review, err = s.repo.Put(ctx, review)
	if err != nil {
		return 0, err
	}

	return review.ID, nil
}

func (s *ReviewService) Update(ctx context.Context, req UpdateRequest) error {
	score, err := domain.NewScore(req.Score)
	if err != nil {
		return err
	}
	review := domain.Review{
		ID:    req.ID,
		Title: req.Title,
		Text:  req.Text,
		Score: score,
	}

	err = s.repo.Update(ctx, review)
	if repoerrors.IsErrNotFound(err) {
		return errors.ErrResourceNotFound{Inner: err}
	}

	return err
}

func (s *ReviewService) Get(ctx context.Context, id int) (domain.Review, error) {
	review, err := s.repo.Get(ctx, id)
	if repoerrors.IsErrNotFound(err) {
		return domain.Review{}, errors.ErrResourceNotFound{Inner: err}
	}

	return review, err
}

func (s *ReviewService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
