package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	reviewpb "github.com/lunn06/library/gateway/internal/api/proto/review"
)

const (
	reviewGetSubj            = "review.get"
	reviewGetAllByBookIdSubj = "review.getAllByBookId"
	reviewPutSubj            = "review.put"
	reviewUpdateSubj         = "review.update"
	reviewDeleteSubj         = "review.delete"
)

func NewReviewAPI(conn *nats.Conn) ReviewAPI {
	return ReviewAPI{conn: conn}
}

type ReviewAPI struct {
	conn *nats.Conn
}

func (ri ReviewAPI) Register(router fiber.Router) {
	router.
		Get("/reviews/book/:id", ri.GetByBookID).
		Name(reviewGetAllByBookIdSubj).
		Get("/review/:id", ri.Get).
		Name(reviewGetSubj).
		Post("/review", ri.Put).
		Name(reviewPutSubj).
		Patch("/review/:id", ri.Update).
		Name(reviewUpdateSubj).
		Delete("/review/:id", ri.Delete).
		Name(reviewDeleteSubj)
}

func (ri ReviewAPI) GetByBookID(ctx *fiber.Ctx) error {
	var params struct {
		bookID int64 `params:"id"`
	}
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	data, err := proto.Marshal(&reviewpb.GetByBookIdRequest{
		BookId: params.bookID,
	})
	if err != nil {
		return err
	}

	respMsg, err := ri.conn.RequestWithContext(ctx.Context(), reviewGetAllByBookIdSubj, data)
	if err != nil {
		return err
	}

	var resp reviewpb.GetByBookIdResponse
	if err = protojson.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ri ReviewAPI) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&reviewpb.GetRequest{ReviewId: int64(id)})
	if err != nil {
		return err
	}

	respMsg, err := ri.conn.RequestWithContext(ctx.Context(), reviewGetSubj, data)
	if err != nil {
		return err
	}

	var resp reviewpb.GetResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ri ReviewAPI) Put(ctx *fiber.Ctx) error {
	var req reviewpb.CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := ri.conn.RequestWithContext(ctx.Context(), reviewPutSubj, data)
	if err != nil {
		return err
	}

	var resp reviewpb.CreateResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ri ReviewAPI) Update(ctx *fiber.Ctx) error {
	var req reviewpb.UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := ri.conn.RequestWithContext(ctx.Context(), reviewUpdateSubj, data)
	if err != nil {
		return err
	}

	var resp reviewpb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ri ReviewAPI) Delete(ctx *fiber.Ctx) error {
	var req reviewpb.DeleteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := ri.conn.RequestWithContext(ctx.Context(), reviewDeleteSubj, data)
	if err != nil {
		return err
	}

	var resp reviewpb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}
