package genre

import (
	"context"

	repoerrors "github.com/lunn06/library/bookinfo/internal/app/repository/errors"
	genrerepo "github.com/lunn06/library/bookinfo/internal/app/repository/genre"
	"github.com/lunn06/library/bookinfo/internal/app/service/errors"
	"github.com/lunn06/library/bookinfo/internal/domain"
)

// TODO: add errors handling

type SearchRequest struct {
	Title  string
	Offset int
	Limit  int
}

type UpdateRequest struct {
	ID          int
	Title       string
	Description string
	BooksIDs    []int
}

type CreateRequest struct {
	Title       string
	Description string
	BooksIDs    []int
}

func NewService(repo genrerepo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo genrerepo.Repo
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (int, error) {
	genre := domain.Genre{
		Title:       req.Title,
		Description: req.Description,
	}
	genre, err := s.repo.Put(ctx, genre, req.BooksIDs...)
	if err != nil {
		return 0, err
	}

	return genre.ID, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) error {
	genre := domain.Genre{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	}

	err := s.repo.Update(ctx, genre, req.BooksIDs...)
	if repoerrors.IsErrNotFound(err) {
		return errors.ErrResourceNotFound{Inner: err}
	}

	return err
}

func (s *Service) Search(ctx context.Context, req SearchRequest) ([]domain.Genre, error) {
	genres, err := s.repo.SearchByTitleWithLimitOffset(ctx, req.Title, req.Limit, req.Offset)
	if repoerrors.IsErrNotFound(err) {
		return nil, errors.ErrResourceNotFound{Inner: err}
	}

	return genres, err
}

func (s *Service) Get(ctx context.Context, id int) (domain.Genre, error) {
	genre, err := s.repo.Get(ctx, id)
	if repoerrors.IsErrNotFound(err) {
		return domain.Genre{}, errors.ErrResourceNotFound{Inner: err}
	}

	return genre, err
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
