package nats

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	authorpb "github.com/lunn06/library/bookinfo/internal/api/proto/author"
	authorservice "github.com/lunn06/library/bookinfo/internal/app/service/author"
	"github.com/lunn06/library/bookinfo/internal/app/service/errors"
	"github.com/lunn06/library/bookinfo/pkg/nats/middleware"
)

func RegisterAuthorConsumer(conn *nats.Conn, cons *AuthorConsumer) error {
	mws := []middleware.Middleware{
		middleware.Recover(),
		middleware.Logger(slog.Default()),
	}
	for subj, handler := range map[string]nats.MsgHandler{
		"author.search": cons.Search,
		"author.get":    cons.Get,
		"author.put":    cons.Put,
		"author.update": cons.Update,
		"author.delete": cons.Delete,
	} {
		_, err := conn.Subscribe(subj, middleware.With(handler, mws...))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewAuthorConsumer(service *authorservice.Service) *AuthorConsumer {
	return &AuthorConsumer{
		service: service,
	}
}

type AuthorConsumer struct {
	service *authorservice.Service
}

func (ac *AuthorConsumer) Search(msg *nats.Msg) {
	var req authorpb.SearchRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.SearchResponse{
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

	search := authorservice.SearchRequest{
		Name:   req.Name,
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	authors, err := ac.service.Search(context.Background(), search)
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
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Get(msg *nats.Msg) {
	var req authorpb.GetRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.GetResponse{
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

	author, err := ac.service.Get(context.Background(), int(req.AuthorId))
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
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Put(msg *nats.Msg) {
	var req authorpb.CreateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.CreateResponse{
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

	create := authorservice.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		BooksIDs:    fromTo[int64, int](req.BooksIds),
	}

	id, err := ac.service.Create(context.Background(), create)
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
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Update(msg *nats.Msg) {
	var req authorpb.UpdateRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.EmptyResponse{
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

	update := authorservice.UpdateRequest{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		BooksIDs:    fromTo[int64, int](req.BooksIds),
	}

	err := ac.service.Update(context.Background(), update)
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
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}

func (ac *AuthorConsumer) Delete(msg *nats.Msg) {
	var req authorpb.DeleteRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		out, err := proto.Marshal(&authorpb.EmptyResponse{
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

	err := ac.service.Delete(context.Background(), int(req.AuthorId))
	if err != nil {
		slog.Error("Error on author delete", "err", err)
		statusCode = http.StatusInternalServerError
	}

	out, err := proto.Marshal(&authorpb.EmptyResponse{
		StatusCode: int32(statusCode),
	})
	if err != nil {
		slog.Error("Error on marshal", "err", err)
		_ = msg.Nak()
		return
	}

	_ = msg.Respond(out)
}
