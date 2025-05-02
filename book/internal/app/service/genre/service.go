package genre

import (
	"context"

	genrerepo "github.com/lunn06/library/book/internal/app/repository/genre"
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

	return s.repo.Update(ctx, genre, req.BooksIDs...)
}

func (s *Service) Search(ctx context.Context, req SearchRequest) ([]domain.Genre, error) {
	return s.repo.SearchByTitleWithLimitOffset(ctx, req.Title, req.Limit, req.Offset)
}

func (s *Service) Get(ctx context.Context, id int) (domain.Genre, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
