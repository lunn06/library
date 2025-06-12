package nats

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	bookpb "github.com/lunn06/library/bookinfo/internal/api/proto/book"
	bookservice "github.com/lunn06/library/bookinfo/internal/app/service/book"
	"github.com/lunn06/library/bookinfo/internal/app/service/errors"
	"github.com/lunn06/library/bookinfo/pkg/nats/middleware"
)

func RegisterBookConsumer(conn *nats.Conn, cons *BookConsumer) error {
	mws := []middleware.Middleware{
		middleware.Recover(),
		middleware.Logger(slog.Default()),
	}
	for subj, handler := range map[string]nats.MsgHandler{
		"book.search": cons.Search,
		"book.get":    cons.Get,
		"book.put":    cons.Put,
		"book.update": cons.Update,
		"book.delete": cons.Delete,
	} {
		_, err := conn.Subscribe(subj, middleware.With(handler, mws...))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewBookConsumer(service *bookservice.Service) *BookConsumer {
	return &BookConsumer{
		service: service,
	}
}

type BookConsumer struct {
	service *bookservice.Service
}

func (bc *BookConsumer) Search(msg *nats.Msg) {
	var req bookpb.SearchRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.SearchResponse{
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

	search := bookservice.SearchRequest{
		Title:  req.Title,
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	books, err := bc.service.Search(context.Background(), search)
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on gateway search", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on gateway search", "err", err)
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
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Get(msg *nats.Msg) {
	var req bookpb.GetRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.GetResponse{
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

	book, err := bc.service.Get(context.Background(), int(req.BookId))
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on gateway get", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on gateway get", "err", err)
		statusCode = http.StatusInternalServerError
	}

	authorsIDs := make([]int64, len(book.Authors))
	for i, author := range book.Authors {
		authorsIDs[i] = int64(author.ID)
	}

	genresIDs := make([]int64, len(book.Genres))
	for i, genre := range book.Genres {
		genresIDs[i] = int64(genre.ID)
	}

	out, err := proto.Marshal(&bookpb.GetResponse{
		Id:          int64(book.ID),
		UserId:      int64(book.UserID),
		Title:       book.Title,
		Description: book.Description,
		BookUrl:     book.BookURL,
		CoverUrl:    book.CoverURL,
		AuthorsIds:  authorsIDs,
		GenresIds:   genresIDs,
		StatusCode:  int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Put(msg *nats.Msg) {
	var req bookpb.CreateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.CreateResponse{
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

	create := bookservice.CreateRequest{
		UserID:      int(req.UserId),
		Title:       req.Title,
		Description: req.Description,
		BookURL:     req.BookUrl,
		CoverURL:    req.CoverUrl,
		AuthorsIDs:  fromTo[int64, int](req.AuthorsIds),
		GenresIDs:   fromTo[int64, int](req.GenresIds),
	}

	id, err := bc.service.Create(context.Background(), create)
	if err != nil {
		slog.Error("Error on gateway put", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.CreateResponse{
		BookId:     int64(id),
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Update(msg *nats.Msg) {
	var req bookpb.UpdateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.EmptyResponse{
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

	err := bc.service.Update(context.Background(), update)
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on gateway update", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on gateway update", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (bc *BookConsumer) Delete(msg *nats.Msg) {
	var req bookpb.DeleteRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&bookpb.EmptyResponse{
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

	err := bc.service.Delete(context.Background(), int(req.BookId))
	if err != nil {
		slog.Error("Error on gateway delete", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&bookpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}
