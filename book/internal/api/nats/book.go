package nats

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	bookpb "github.com/lunn06/library/book/internal/api/proto/book"
	bookservice "github.com/lunn06/library/book/internal/app/service/book"
)

func RegisterBookConsumer(conn *nats.Conn, cons *BookConsumer) error {
	for subj, handler := range map[string]nats.MsgHandler{
		"book.search": cons.Search,
		"book.get":    cons.Get,
		"book.put":    cons.Put,
		"book.update": cons.Update,
		"book.delete": cons.Delete,
	} {
		_, err := conn.Subscribe(subj, handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewBookConsumer(cfg Config, service *bookservice.Service) *BookConsumer {
	return &BookConsumer{
		service: service,
		timeout: cfg.RequestTimeout,
	}
}

type BookConsumer struct {
	service *bookservice.Service
	timeout time.Duration
}

func (bc *BookConsumer) Search(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req bookpb.SearchRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.SearchResponse{
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

	search := bookservice.SearchRequest{
		Title:  req.Title,
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()

	books, err := bc.service.Search(ctx, search)
	if err != nil {
		slog.Error("Error on book search", "err", err)
		statusCode = http.StatusInternalServerError
	}

	items := make([]*bookpb.SearchItem, len(books))
	for i, book := range books {
		items[i] = &bookpb.SearchItem{
			Id:          int64(book.ID),
			UserId:      int64(book.UserID),
			Title:       book.Title,
			Description: book.Description,
			BookUrl:     book.BookURL,
			CoverUrl:    book.CoverURL,
		}
	}

	out, err := proto.Marshal(&bookpb.SearchResponse{
		Items:      items,
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Get(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req bookpb.GetRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.GetResponse{
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

	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()

	book, err := bc.service.Get(ctx, int(req.BookId))
	if err != nil {
		slog.Error("Error on book get", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.GetResponse{
		Id:          int64(book.ID),
		UserId:      int64(book.UserID),
		Title:       book.Title,
		Description: book.Description,
		BookUrl:     book.BookURL,
		CoverUrl:    book.CoverURL,
		StatusCode:  int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Put(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req bookpb.CreateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.CreateResponse{
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

	create := bookservice.CreateRequest{
		UserID:      int(req.UserId),
		Title:       req.Title,
		Description: req.Description,
		BookURL:     req.BookUrl,
		CoverURL:    req.CoverUrl,
		AuthorsIDs:  fromTo[int64, int](req.AuthorsIds),
		GenresIDs:   fromTo[int64, int](req.GenresIds),
	}

	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()

	id, err := bc.service.Create(ctx, create)
	if err != nil {
		slog.Error("Error on book put", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.CreateResponse{
		BookId:     int64(id),
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Update(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req bookpb.UpdateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.EmptyResponse{
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

	update := bookservice.UpdateRequest{
		ID:          int(req.Id),
		UserID:      int(req.UserId),
		Title:       req.Title,
		Description: req.Description,
		BookURL:     req.BookUrl,
		CoverURL:    req.CoverUrl,
		AuthorsIDs:  fromTo[int64, int](req.AuthorsIds),
		GenresIDs:   fromTo[int64, int](req.GenresIds),
	}

	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()

	err := bc.service.Update(ctx, update)
	if err != nil {
		slog.Error("Error on book update", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Delete(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req bookpb.DeleteRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.EmptyResponse{
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

	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()

	err := bc.service.Delete(ctx, int(req.BookId))
	if err != nil {
		slog.Error("Error on book delete", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}
