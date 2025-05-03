package nats

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	authorpb "github.com/lunn06/library/book/internal/api/proto/author"
	authorservice "github.com/lunn06/library/book/internal/app/service/author"
	"github.com/lunn06/library/book/internal/app/service/errors"
)

func RegisterAuthorConsumer(conn *nats.Conn, cons *AuthorConsumer) error {
	for subj, handler := range map[string]nats.MsgHandler{
		"author.search": cons.Search,
		"author.get":    cons.Get,
		"author.put":    cons.Put,
		"author.update": cons.Update,
		"author.delete": cons.Delete,
	} {
		_, err := conn.Subscribe(subj, handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewAuthorConsumer(cfg Config, service *authorservice.Service) *AuthorConsumer {
	return &AuthorConsumer{
		service: service,
		timeout: cfg.RequestTimeout,
	}
}

type AuthorConsumer struct {
	service *authorservice.Service
	timeout time.Duration
}

func (ac *AuthorConsumer) Search(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req authorpb.SearchRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.SearchResponse{
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

	search := authorservice.SearchRequest{
		Name:   req.Name,
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	ctx, cancel := context.WithTimeout(context.Background(), ac.timeout)
	defer cancel()

	authors, err := ac.service.Search(ctx, search)
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on author search", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on author search", "err", err)
		statusCode = http.StatusInternalServerError
	}

	items := make([]*authorpb.SearchItem, len(authors))
	for i, author := range authors {
		items[i] = &authorpb.SearchItem{
			Id:          int64(author.ID),
			Name:        author.Name,
			Description: author.Description,
		}
	}

	out, err := proto.Marshal(&authorpb.SearchResponse{
		Items:      items,
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Get(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req authorpb.GetRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.GetResponse{
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

	ctx, cancel := context.WithTimeout(context.Background(), ac.timeout)
	defer cancel()

	author, err := ac.service.Get(ctx, int(req.AuthorId))
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on author get", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on author get", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&authorpb.GetResponse{
		Id:          int64(author.ID),
		Name:        author.Name,
		Description: author.Description,
		StatusCode:  int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Put(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req authorpb.CreateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.CreateResponse{
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

	create := authorservice.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		BooksIDs:    fromTo[int64, int](req.BooksIds),
	}

	ctx, cancel := context.WithTimeout(context.Background(), ac.timeout)
	defer cancel()

	id, err := ac.service.Create(ctx, create)
	if err != nil {
		slog.Error("Error on author put", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&authorpb.CreateResponse{
		AuthorId:   int64(id),
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Update(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req authorpb.UpdateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.EmptyResponse{
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

	update := authorservice.UpdateRequest{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		BooksIDs:    fromTo[int64, int](req.BooksIds),
	}

	ctx, cancel := context.WithTimeout(context.Background(), ac.timeout)
	defer cancel()

	err := ac.service.Update(ctx, update)
	if errors.IsErrResourceNotFound(err) {
		slog.Error("Not found on author update", "err", err)
		statusCode = http.StatusNotFound
	} else if err != nil {
		slog.Error("Error on author update", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&authorpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Delete(msg *nats.Msg) {
	slog.Info("Data", "msg", msg.Data)

	var req authorpb.DeleteRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.EmptyResponse{
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

	ctx, cancel := context.WithTimeout(context.Background(), ac.timeout)
	defer cancel()

	err := ac.service.Delete(ctx, int(req.AuthorId))
	if err != nil {
		slog.Error("Error on author delete", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&authorpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		return
	}

	_ = msg.Respond(out)
}
