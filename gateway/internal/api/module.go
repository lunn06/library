package api

import (
	"github.com/lunn06/library/gateway/internal/api/server"
	"go.uber.org/fx"

	"github.com/lunn06/library/gateway/internal/api/nats"
)

var Module = fx.Module("api",
	nats.Module,
	server.Module,

	fx.Provide(
		NewBookInfoAPI,
		NewAuthorAPI,
		NewGenreAPI,
		NewBookFileAPI,
		NewReviewAPI,

		NewBookFileClient,
	),

	fx.Invoke(
		registerRouters,
	),
)

func registerRouters(
	server *server.Server,
	bookInfo BookInfoAPI,
	author AuthorAPI,
	genre GenreAPI,
	bookFile BookFileAPI,
	review ReviewAPI,
) {
	router := server.Router()

	bookInfo.Register(router)
	author.Register(router)
	genre.Register(router)
	bookFile.Register(router)
	review.Register(router)
}
