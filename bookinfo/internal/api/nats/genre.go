package nats

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	genrepb "github.com/lunn06/library/bookinfo/internal/api/proto/genre"
	"github.com/lunn06/library/bookinfo/internal/app/service/errors"
	genreservice "github.com/lunn06/library/bookinfo/internal/app/service/genre"
	"github.com/lunn06/library/bookinfo/pkg/nats/middleware"
)

func RegisterGenreConsumer(conn *nats.Conn, cons *GenreConsumer) error {
	mws := []middleware.Middleware{
		middleware.Recover(),
		middleware.Logger(slog.Default()),
	}
	for subj, handler := range map[string]nats.MsgHandler{
		"genre.search": cons.Search,
		"genre.get":    cons.Get,
		"genre.put":    cons.Put,
		"genre.update": cons.Update,
		"genre.delete": cons.Delete,
	} {
		_, err := conn.Subscribe(subj, middleware.With(handler, mws...))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGenreConsumer(service *genreservice.Service) *GenreConsumer {
	return &GenreConsumer{
		service: service,
	}
}

type GenreConsumer struct {
	service *genreservice.Service
}

func (gc *GenreConsumer) Search(msg *nats.Msg) {
	var req genrepb.SearchRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.SearchResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			_ = msg.Nak()
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	search := genreservice.SearchRequest{
		Title:  req.Title,
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	books, err := gc.service.Search(context.Background(), search)
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on genre search", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on genre search", "err", err)
		statusCode = http.StatusInternalServerError
	}

	items := make([]*genrepb.SearchItem, len(books))
	for i, book := range books {
		items[i] = &genrepb.SearchItem{
			Id:          int64(book.ID),
			Title:       book.Title,
			Description: book.Description,
		}
	}

	out, err := proto.Marshal(&genrepb.SearchResponse{
		Items:      items,
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Get(msg *nats.Msg) {
	var req genrepb.GetRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.GetResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			_ = msg.Nak()
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	genre, err := gc.service.Get(context.Background(), int(req.GenreId))
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on genre get", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on genre get", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.GetResponse{
		Id:          int64(genre.ID),
		Title:       genre.Title,
		Description: genre.Description,
		StatusCode:  int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Put(msg *nats.Msg) {
	var req genrepb.CreateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.CreateResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			_ = msg.Nak()
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	create := genreservice.CreateRequest{
		Title:       req.Title,
		Description: req.Description,
		BooksIDs:    fromTo[int64, int](req.BooksIds),
	}

	id, err := gc.service.Create(context.Background(), create)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.CreateResponse{
		GenreId:    int64(id),
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Update(msg *nats.Msg) {
	var req genrepb.UpdateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.EmptyResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			_ = msg.Nak()
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	update := genreservice.UpdateRequest{
		ID:          int(req.Id),
		Title:       req.Title,
		Description: req.Description,
		BooksIDs:    fromTo[int64, int](req.BooksIds),
	}

	err := gc.service.Update(context.Background(), update)
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on genre update", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on genre update", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Delete(msg *nats.Msg) {
	var req genrepb.DeleteRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.EmptyResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			_ = msg.Nak()
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	err := gc.service.Delete(context.Background(), int(req.GenreId))
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}
