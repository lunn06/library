package nats

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	genrepb "github.com/lunn06/library/book/internal/api/proto/genre"
	genreservice "github.com/lunn06/library/book/internal/app/service/genre"
)

func RegisterGenreConsumer(conn *nats.Conn, cons *GenreConsumer) error {
	for subj, handler := range map[string]nats.MsgHandler{
		"genre.search": cons.Search,
		"genre.get":    cons.Get,
		"genre.put":    cons.Put,
		"genre.update": cons.Update,
		"genre.delete": cons.Delete,
	} {
		_, err := conn.Subscribe(subj, handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGenreConsumer(cfg Config, service *genreservice.Service) *GenreConsumer {
	return &GenreConsumer{
		service: service,
		timeout: cfg.RequestTimeout,
	}
}

type GenreConsumer struct {
	service *genreservice.Service
	timeout time.Duration
}

func (gc *GenreConsumer) Search(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req genrepb.SearchRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.SearchResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
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

	ctx, cancel := context.WithTimeout(context.Background(), gc.timeout)
	defer cancel()

	books, err := gc.service.Search(ctx, search)
	if err != nil {
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
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Get(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req genrepb.GetRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.GetResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	ctx, cancel := context.WithTimeout(context.Background(), gc.timeout)
	defer cancel()

	genre, err := gc.service.Get(ctx, int(req.GenreId))
	if err != nil {
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
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Put(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req genrepb.CreateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.CreateResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
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

	ctx, cancel := context.WithTimeout(context.Background(), gc.timeout)
	defer cancel()

	id, err := gc.service.Create(ctx, create)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.CreateResponse{
		GenreId:    int64(id),
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Update(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req genrepb.UpdateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.EmptyResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
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

	ctx, cancel := context.WithTimeout(context.Background(), gc.timeout)
	defer cancel()

	err := gc.service.Update(ctx, update)
	if err != nil {
		slog.Error("Error on update", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (gc *GenreConsumer) Delete(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req genrepb.DeleteRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&genrepb.EmptyResponse{
			StatusCode: http.StatusUnprocessableEntity,
		})
		if err != nil {
			slog.Error("Error on marshal", "err", err)
			return
		}
		_ = msg.Respond(out)
		return
	}

	statusCode := http.StatusOK

	ctx, cancel := context.WithTimeout(context.Background(), gc.timeout)
	defer cancel()

	err := gc.service.Delete(ctx, int(req.GenreId))
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&genrepb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}
