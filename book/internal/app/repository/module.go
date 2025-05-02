package repository

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/book/internal/app/repository/author"
	"github.com/lunn06/library/book/internal/app/repository/book"
	"github.com/lunn06/library/book/internal/app/repository/genre"
)

var Module = fx.Options(
	author.Module,
	book.Module,
	genre.Module,
)
