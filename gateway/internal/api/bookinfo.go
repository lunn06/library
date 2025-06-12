package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	bookpb "github.com/lunn06/library/gateway/internal/api/proto/book"
)

const (
	bookInfoSearchSubj = "book.search"
	bookInfoGetSubj    = "book.get"
	bookInfoPutSubj    = "book.put"
	bookInfoUpdateSubj = "book.update"
	bookInfoDeleteSubj = "book.delete"
)

func NewBookInfoAPI(conn *nats.Conn) BookInfoAPI {
	return BookInfoAPI{conn: conn}
}

type BookInfoAPI struct {
	conn *nats.Conn
}

func (bi BookInfoAPI) Register(router fiber.Router) {
	router.
		Get("/book/search/:title", bi.Search).
		Name(bookInfoSearchSubj).
		Get("/book/:id", bi.Get).
		Name(bookInfoGetSubj).
		Post("/book", bi.Put).
		Name(bookInfoPutSubj).
		Patch("/book/:id", bi.Update).
		Name(bookInfoUpdateSubj).
		Delete("/book/:id", bi.Delete).
		Name(bookInfoDeleteSubj)
}

func (bi BookInfoAPI) Search(ctx *fiber.Ctx) error {
	var params struct {
		offset int32 `query:"offset"`
		limit  int32 `query:"limit"`
	}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}

	data, err := proto.Marshal(&bookpb.SearchRequest{
		Title:  ctx.Params("title"),
		Offset: params.offset,
		Limit:  params.limit,
	})
	if err != nil {
		return err
	}

	respMsg, err := bi.conn.RequestWithContext(ctx.Context(), bookInfoSearchSubj, data)
	if err != nil {
		return err
	}

	var resp bookpb.SearchResponse
	if err = protojson.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (bi BookInfoAPI) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&bookpb.GetRequest{BookId: int64(id)})
	if err != nil {
		return err
	}

	respMsg, err := bi.conn.RequestWithContext(ctx.Context(), bookInfoGetSubj, data)
	if err != nil {
		return err
	}

	var resp bookpb.GetResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (bi BookInfoAPI) Put(ctx *fiber.Ctx) error {
	var req bookpb.CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := bi.conn.RequestWithContext(ctx.Context(), bookInfoPutSubj, data)
	if err != nil {
		return err
	}

	var resp bookpb.CreateResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (bi BookInfoAPI) Update(ctx *fiber.Ctx) error {
	var req bookpb.UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := bi.conn.RequestWithContext(ctx.Context(), bookInfoUpdateSubj, data)
	if err != nil {
		return err
	}

	var resp bookpb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (bi BookInfoAPI) Delete(ctx *fiber.Ctx) error {
	var req bookpb.DeleteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := bi.conn.RequestWithContext(ctx.Context(), bookInfoDeleteSubj, data)
	if err != nil {
		return err
	}

	var resp bookpb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}
