package author

import (
	"context"

	authorepo "github.com/lunn06/library/book/internal/app/repository/author"
	"github.com/lunn06/library/book/internal/domain"
)

// TODO: add errors handling

type SearchRequest struct {
	Name   string
	Offset int
	Limit  int
}

type UpdateRequest struct {
	ID          int
	Name        string
	Description string
	BooksIDs    []int
}

type CreateRequest struct {
	Name        string
	Description string
	BooksIDs    []int
}

func NewService(repo authorepo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo authorepo.Repo
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (int, error) {
	author := domain.Author{
		Name:        req.Name,
		Description: req.Description,
	}
	author, err := s.repo.Put(ctx, author, req.BooksIDs...)
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) error {
	author := domain.Author{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	return s.repo.Update(ctx, author, req.BooksIDs...)
}

func (s *Service) Search(ctx context.Context, req SearchRequest) ([]domain.Author, error) {
	return s.repo.SearchByNameWithLimitOffset(ctx, req.Name, req.Limit, req.Offset)
}

func (s *Service) Get(ctx context.Context, id int) (domain.Author, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
