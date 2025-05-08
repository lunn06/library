package client

import (
	"github.com/nats-io/nats.go/jetstream"

	"github.com/lunn06/library/bookfile/internal/app/repository"
	"github.com/lunn06/library/bookfile/internal/app/service"
	"github.com/lunn06/library/bookfile/internal/infrastructure/db/js"
)

func New(config Config) (*Client, error) {
	jstream, err := js.NewJetStream(config.URL)
	if err != nil {
		return nil, err
	}

	objStore, err := js.NewObjectStore(jstream)
	if err != nil {
		return nil, err
	}

	repo := repository.NewJsBookRepo(objStore)
	bookService := service.NewBookService(repo)

	return &Client{
		BookService: bookService,
		jstream:     jstream,
	}, nil
}

type Client struct {
	*service.BookService
	jstream jetstream.JetStream
}

func (c *Client) Close() error {
	return c.jstream.Conn().Drain()
}
