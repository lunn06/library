//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	natscontainer "github.com/testcontainers/testcontainers-go/modules/nats"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/protobuf/proto"

	natsapi "github.com/lunn06/library/bookinfo/internal/api/nats"
	authorepo "github.com/lunn06/library/bookinfo/internal/app/repository/author"
	bookrepo "github.com/lunn06/library/bookinfo/internal/app/repository/book"
	genrerepo "github.com/lunn06/library/bookinfo/internal/app/repository/genre"
	"github.com/lunn06/library/bookinfo/internal/app/service/author"
	"github.com/lunn06/library/bookinfo/internal/app/service/book"
	"github.com/lunn06/library/bookinfo/internal/app/service/genre"
	"github.com/lunn06/library/bookinfo/internal/config"
	"github.com/lunn06/library/bookinfo/internal/infrastructure/db/postgres"
)

const reqTimeout = time.Second

const (
	authorSearchSubj = "author.search"
	authorGetSubj    = "author.get"
	authorPutSubj    = "author.put"
	authorUpdateSubj = "author.update"
	authorDeleteSubj = "author.delete"

	bookSearchSubj = "review.search"
	bookGetSubj    = "review.get"
	bookPutSubj    = "review.put"
	bookUpdateSubj = "review.update"
	bookDeleteSubj = "review.delete"

	genreSearchSubj = "genre.search"
	genreGetSubj    = "genre.get"
	genrePutSubj    = "genre.put"
	genreUpdateSubj = "genre.update"
	genreDeleteSubj = "genre.delete"
)

const (
	postgresUser     = "test-review-user"
	postgresPassword = "test-review-password"
	postgresDb       = "test-review-db"
	postgresSslMode  = "disable"
)

var nc *nats.Conn

func TestMain(m *testing.M) {
	natsC := natsContainer()
	natsEndpoint, err := natsC.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	postgresC := postgresContainer()
	postgresEndpoint, err := postgresC.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	cfg.Nats = natsapi.Config{
		URL: fmt.Sprintf("nats://%s", natsEndpoint),
	}
	cfg.Postgres.URL = postgresEndpoint
	cfg.Postgres.User = postgresUser
	cfg.Postgres.Password = postgresPassword
	cfg.Postgres.DB = postgresDb
	cfg.Postgres.SslMode = postgresSslMode

	entClient, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	authorRepo := authorepo.NewEntRepo(entClient)
	authorService := author.NewService(authorRepo)
	authorCons := natsapi.NewAuthorConsumer(authorService)

	bookRepo := bookrepo.NewEntRepo(entClient)
	bookService := book.NewService(bookRepo)
	bookCons := natsapi.NewBookConsumer(bookService)

	genreRepo := genrerepo.NewEntRepo(entClient)
	genreService := genre.NewService(genreRepo)
	genreCons := natsapi.NewGenreConsumer(genreService)

	nc, err = nats.Connect(cfg.Nats.URL)

	if err = natsapi.RegisterAuthorConsumer(nc, authorCons); err != nil {
		panic(err)
	}
	if err = natsapi.RegisterBookConsumer(nc, bookCons); err != nil {
		panic(err)
	}
	if err = natsapi.RegisterGenreConsumer(nc, genreCons); err != nil {
		panic(err)
	}

	m.Run()

	err = errors.Join(
		nc.Drain(),
		natsC.Terminate(context.Background()),
		postgresC.Terminate(context.Background()),
	)
	if err != nil {
		panic(err)
	}
}

func natsContainer() testcontainers.Container {
	ctx := context.Background()
	natsC, err := natscontainer.Run(ctx, "nats:2.11-alpine3.21")
	if err != nil {
		panic(err)
	}

	return natsC
}

func postgresContainer() testcontainers.Container {
	ctx := context.Background()
	postgresC, err := postgrescontainer.Run(ctx,
		"postgres:17-alpine3.21",
		postgrescontainer.WithUsername(postgresUser),
		postgrescontainer.WithPassword(postgresPassword),
		postgrescontainer.WithDatabase(postgresDb),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		panic(err)
	}

	return postgresC
}

func request(subj string, req, resp proto.Message) error {
	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	resMsg, err := nc.Request(subj, data, reqTimeout)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(resMsg.Data, resp)
	if err != nil {
		return err
	}

	return nil
}
