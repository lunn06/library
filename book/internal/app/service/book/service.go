package book

import (
	"context"

	bookrepo "github.com/lunn06/library/book/internal/app/repository/book"
	repoerrors "github.com/lunn06/library/book/internal/app/repository/errors"
	"github.com/lunn06/library/book/internal/app/service/errors"
	"github.com/lunn06/library/book/internal/domain"
)

// TODO: add errors handling

type SearchRequest struct {
	Title  string
	Offset int
	Limit  int
}

type UpdateRequest struct {
	ID          int
	UserID      int
	Title       string
	Description string
	BookURL     string
	CoverURL    *string
	AuthorsIDs  []int
	GenresIDs   []int
}

type CreateRequest struct {
	UserID      int
	Title       string
	Description string
	BookURL     string
	CoverURL    *string
	AuthorsIDs  []int
	GenresIDs   []int
}

func NewService(repo bookrepo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo bookrepo.Repo
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (int, error) {
	book := domain.Book{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		BookURL:     req.BookURL,
		CoverURL:    req.CoverURL,
	}
	book, err := s.repo.Put(ctx, book, req.AuthorsIDs, req.GenresIDs)
	if err != nil {
		return 0, err
	}

	return book.ID, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) error {
	book := domain.Book{
		ID:          req.ID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		BookURL:     req.BookURL,
		CoverURL:    req.CoverURL,
	}

	err := s.repo.Update(ctx, book, req.AuthorsIDs, req.GenresIDs)
	if repoerrors.IsErrNotFound(err) {
		return errors.ErrResourceNotFound{Inner: err}
	}

	return err
}

func (s *Service) Search(ctx context.Context, req SearchRequest) ([]domain.Book, error) {
	books, err := s.repo.SearchByTitleWithLimitOffset(ctx, req.Title, req.Limit, req.Offset)
	if repoerrors.IsErrNotFound(err) {
		return nil, errors.ErrResourceNotFound{Inner: err}
	}

	return books, err
}

func (s *Service) Get(ctx context.Context, id int) (domain.Book, error) {
	book, err := s.repo.Get(ctx, id)
	if repoerrors.IsErrNotFound(err) {
		return domain.Book{}, errors.ErrResourceNotFound{Inner: err}
	}

	return book, err
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
