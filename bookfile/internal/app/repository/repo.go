package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/lunn06/library/bookfile/internal/domain"
)

type BookRepo interface {
	Get(ctx context.Context, bookUUID uuid.UUID) (domain.Book, error)
	Put(ctx context.Context, book domain.Book) error
	Delete(ctx context.Context, bookUUID uuid.UUID) error
}
