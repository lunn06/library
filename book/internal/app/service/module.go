package service

import (
	"go.uber.org/fx"

	"github.com/lunn06/library/book/internal/app/service/author"
	"github.com/lunn06/library/book/internal/app/service/book"
	"github.com/lunn06/library/book/internal/app/service/genre"
)

var Module = fx.Options(
	author.Module,
	book.Module,
	genre.Module,
)
