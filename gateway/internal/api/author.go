package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	authorpb "github.com/lunn06/library/gateway/internal/api/proto/author"
)

const (
	authorSearchSubj = "author.search"
	authorGetSubj    = "author.get"
	authorPutSubj    = "author.put"
	authorUpdateSubj = "author.update"
	authorDeleteSubj = "author.delete"
)

func NewAuthorAPI(conn *nats.Conn) AuthorAPI {
	return AuthorAPI{conn: conn}
}

type AuthorAPI struct {
	conn *nats.Conn
}

func (ai AuthorAPI) Register(router fiber.Router) {
	router.
		Get("/author/search/:name", ai.Search).
		Name(authorSearchSubj).
		Get("/author/:id", ai.Get).
		Name(authorGetSubj).
		Post("/author", ai.Put).
		Name(authorPutSubj).
		Patch("/author/:id", ai.Update).
		Name(authorUpdateSubj).
		Delete("/author/:id", ai.Delete).
		Name(authorDeleteSubj)
}

func (ai AuthorAPI) Search(ctx *fiber.Ctx) error {
	var params struct {
		offset int32 `query:"offset"`
		limit  int32 `query:"limit"`
	}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}

	data, err := proto.Marshal(&authorpb.SearchRequest{
		Name:   ctx.Params("name"),
		Offset: params.offset,
		Limit:  params.limit,
	})
	if err != nil {
		return err
	}

	respMsg, err := ai.conn.RequestWithContext(ctx.Context(), authorSearchSubj, data)
	if err != nil {
		return err
	}

	var resp authorpb.SearchResponse
	if err = protojson.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ai AuthorAPI) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&authorpb.GetRequest{AuthorId: int64(id)})
	if err != nil {
		return err
	}

	respMsg, err := ai.conn.RequestWithContext(ctx.Context(), authorGetSubj, data)
	if err != nil {
		return err
	}

	var resp authorpb.GetResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ai AuthorAPI) Put(ctx *fiber.Ctx) error {
	var req authorpb.CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := ai.conn.RequestWithContext(ctx.Context(), authorPutSubj, data)
	if err != nil {
		return err
	}

	var resp authorpb.CreateResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ai AuthorAPI) Update(ctx *fiber.Ctx) error {
	var req authorpb.UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := ai.conn.RequestWithContext(ctx.Context(), authorUpdateSubj, data)
	if err != nil {
		return err
	}

	var resp authorpb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (ai AuthorAPI) Delete(ctx *fiber.Ctx) error {
	var req authorpb.DeleteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := ai.conn.RequestWithContext(ctx.Context(), authorDeleteSubj, data)
	if err != nil {
		return err
	}

	var resp authorpb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}
