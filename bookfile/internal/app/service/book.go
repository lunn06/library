package service

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/google/uuid"

	"github.com/lunn06/library/bookfile/internal/app/repository"
	repoerrors "github.com/lunn06/library/bookfile/internal/app/repository/errors"
	servicerrors "github.com/lunn06/library/bookfile/internal/app/service/errors"
	"github.com/lunn06/library/bookfile/internal/domain"
)

// TODO: add servicerrors handling

func NewBookService(repo repository.BookRepo) *BookService {
	return &BookService{
		repo: repo,
	}
}

type BookService struct {
	repo repository.BookRepo
}

func (s *BookService) Create(ctx context.Context, fileName string, r io.ReadCloser) (uuid.UUID, error) {
	switch {
	case fileName == "":
		fileName = "unknown"
	case !strings.HasSuffix(fileName, ".pdf"):
		return uuid.Nil, errors.New("invalid file name")
	}
	if r == nil {
		return uuid.Nil, errors.New("buffer is nil")
	}

	bookUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	book := domain.Book{
		UUID:           bookUUID,
		FileName:       fileName,
		FileReadCloser: r,
	}
	if err = s.repo.Put(ctx, book); err != nil {
		return uuid.Nil, err
	}

	return bookUUID, nil
}

func (s *BookService) Get(ctx context.Context, bookUUID uuid.UUID) (domain.Book, error) {
	book, err := s.repo.Get(ctx, bookUUID)
	if repoerrors.IsErrNotFound(err) {
		return domain.Book{}, servicerrors.ErrResourceNotFound{Inner: err}
	}
	if book.FileName == "" {
		book.FileName = "unknown"
	}

	return book, err
}

func (s *BookService) Delete(ctx context.Context, bookUUID uuid.UUID) error {
	return s.repo.Delete(ctx, bookUUID)
}
