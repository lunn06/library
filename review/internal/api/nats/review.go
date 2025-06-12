package nats

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/lunn06/library/bookinfo/pkg/nats/middleware"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	reviewpb "github.com/lunn06/library/review/internal/api/proto/review"
	"github.com/lunn06/library/review/internal/app/service"
	"github.com/lunn06/library/review/internal/app/service/errors"
)

func RegisterReviewConsumer(conn *nats.Conn, cons *ReviewConsumer) error {
	mws := []middleware.Middleware{
		middleware.Recover(),
		middleware.Logger(slog.Default()),
	}
	for subj, handler := range map[string]nats.MsgHandler{
		"review.get":            cons.Get,
		"review.getAllByBookId": cons.GetAllByBookID,
		"review.put":            cons.Put,
		"review.update":         cons.Update,
		"review.delete":         cons.Delete,
	} {
		_, err := conn.Subscribe(subj, middleware.With(handler, mws...))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewReviewConsumer(service *service.ReviewService) *ReviewConsumer {
	return &ReviewConsumer{
		service: service,
	}
}

type ReviewConsumer struct {
	service *service.ReviewService
}

func (rc ReviewConsumer) GetAllByBookID(msg *nats.Msg) {
	var (
		req  reviewpb.GetByBookIdRequest
		resp reviewpb.GetByBookIdResponse
		err  error

		responder = NewResponder(msg)
	)
	defer func() {
		if err != nil {
			slog.Error("Error on reviews get", "err", err)
		}
		responder.Respond(&resp)
	}()
	if err = proto.Unmarshal(msg.Data, &req); err != nil {
		resp.StatusCode = int32(http.StatusUnprocessableEntity)
		return
	}

	reviews, err := rc.service.GetAllByBookID(context.Background(), int(req.BookId))
	if errors.IsErrResourceNotFound(err) {
		resp.StatusCode = http.StatusNotFound
		return
	} else if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		return
	}

	reviewItems := make([]*reviewpb.ReviewItem, len(reviews))
	for i, review := range reviews {
		reviewItems[i] = &reviewpb.ReviewItem{
			Id:        int64(review.ID),
			UserId:    int64(review.UserID),
			BookId:    int64(review.BookID),
			CreatedAt: review.CreatedAt.Unix(),
			Title:     review.Title,
			Text:      review.Text,
			Score:     int32(review.Score),
		}
	}

	resp = reviewpb.GetByBookIdResponse{
		Reviews:    reviewItems,
		StatusCode: http.StatusOK,
	}
}

func (rc ReviewConsumer) Get(msg *nats.Msg) {
	var (
		req  reviewpb.GetRequest
		resp reviewpb.GetResponse
		err  error

		responder = NewResponder(msg)
	)
	defer func() {
		if err != nil {
			slog.Error("Error on review get", "err", err)
		}
		responder.Respond(&resp)
	}()
	if err = proto.Unmarshal(msg.Data, &req); err != nil {
		resp.StatusCode = int32(http.StatusUnprocessableEntity)
		return
	}

	review, err := rc.service.Get(context.Background(), int(req.ReviewId))
	if errors.IsErrResourceNotFound(err) {
		resp.StatusCode = http.StatusNotFound
		return
	} else if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		return
	}

	resp = reviewpb.GetResponse{
		Id:         int64(review.ID),
		UserId:     int64(review.UserID),
		BookId:     int64(review.BookID),
		CreatedAt:  review.CreatedAt.Unix(),
		Title:      review.Title,
		Text:       review.Text,
		Score:      int32(review.Score),
		StatusCode: int32(http.StatusOK),
	}
}

func (rc ReviewConsumer) Put(msg *nats.Msg) {
	var (
		req  reviewpb.CreateRequest
		resp reviewpb.CreateResponse
		err  error

		responder = NewResponder(msg)
	)
	defer func() {
		if err != nil {
			slog.Error("Error on review put", "err", err)
		}
		responder.Respond(&resp)
	}()
	if err = proto.Unmarshal(msg.Data, &req); err != nil {
		resp.StatusCode = int32(http.StatusUnprocessableEntity)
		return
	}

	if err = proto.Unmarshal(msg.Data, &req); err != nil {
		resp.StatusCode = int32(http.StatusUnprocessableEntity)
		return
	}

	create := service.CreateRequest{
		UserID: int(req.UserId),
		BookID: int(req.BookId),
		Title:  req.Title,
		Text:   req.Text,
		Score:  int(req.Score),
	}

	id, err := rc.service.Create(context.Background(), create)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
	}

	resp = reviewpb.CreateResponse{
		ReviewId:   int64(id),
		StatusCode: int32(http.StatusOK),
	}
}

func (rc ReviewConsumer) Update(msg *nats.Msg) {
	var (
		req  reviewpb.UpdateRequest
		resp reviewpb.EmptyResponse
		err  error

		responder = NewResponder(msg)
	)
	defer func() {
		if err != nil {
			slog.Error("Error on review update", "err", err)
		}
		responder.Respond(&resp)
	}()
	if err = proto.Unmarshal(msg.Data, &req); err != nil {
		resp.StatusCode = int32(http.StatusUnprocessableEntity)
		return
	}

	update := service.UpdateRequest{
		ID:    int(req.Id),
		Title: req.Title,
		Text:  req.Text,
		Score: int(req.Score),
	}

	err = rc.service.Update(context.Background(), update)
	if errors.IsErrResourceNotFound(err) {
		resp.StatusCode = http.StatusNotFound
	} else if err != nil {
		resp.StatusCode = http.StatusInternalServerError
	}

	resp.StatusCode = int32(http.StatusOK)
}

func (rc ReviewConsumer) Delete(msg *nats.Msg) {
	var (
		req  reviewpb.DeleteRequest
		resp reviewpb.EmptyResponse
		err  error

		responder = NewResponder(msg)
	)
	defer func() {
		if err != nil {
			slog.Error("Error on review delete", "err", err)
		}
		responder.Respond(&resp)
	}()
	if err = proto.Unmarshal(msg.Data, &req); err != nil {
		resp.StatusCode = http.StatusUnprocessableEntity
		return
	}

	err = rc.service.Delete(context.Background(), int(req.ReviewId))
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		return
	}

	resp.StatusCode = int32(http.StatusOK)
}
