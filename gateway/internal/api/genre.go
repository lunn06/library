package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	genrepb "github.com/lunn06/library/gateway/internal/api/proto/genre"
)

const (
	genreSearchSubj = "genre.search"
	genreGetSubj    = "genre.get"
	genrePutSubj    = "genre.put"
	genreUpdateSubj = "genre.update"
	genreDeleteSubj = "genre.delete"
)

func NewGenreAPI(conn *nats.Conn) GenreAPI {
	return GenreAPI{conn: conn}
}

type GenreAPI struct {
	conn *nats.Conn
}

func (gi GenreAPI) Register(router fiber.Router) {
	router.
		Get("/genre/search/:title", gi.Search).
		Name(genreSearchSubj).
		Get("/genre/:id", gi.Get).
		Name(genreGetSubj).
		Post("/genre", gi.Put).
		Name(genrePutSubj).
		Patch("/genre/:id", gi.Update).
		Name(genreUpdateSubj).
		Delete("/genre/:id", gi.Delete).
		Name(genreDeleteSubj)
}

func (gi GenreAPI) Search(ctx *fiber.Ctx) error {
	var params struct {
		offset int32 `query:"offset"`
		limit  int32 `query:"limit"`
	}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}

	data, err := proto.Marshal(&genrepb.SearchRequest{
		Title:  ctx.Params("title"),
		Offset: params.offset,
		Limit:  params.limit,
	})
	if err != nil {
		return err
	}

	respMsg, err := gi.conn.RequestWithContext(ctx.Context(), genreSearchSubj, data)
	if err != nil {
		return err
	}

	var resp genrepb.SearchResponse
	if err = protojson.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (gi GenreAPI) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&genrepb.GetRequest{GenreId: int64(id)})
	if err != nil {
		return err
	}

	respMsg, err := gi.conn.RequestWithContext(ctx.Context(), genreGetSubj, data)
	if err != nil {
		return err
	}

	var resp genrepb.GetResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (gi GenreAPI) Put(ctx *fiber.Ctx) error {
	var req genrepb.CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := gi.conn.RequestWithContext(ctx.Context(), genrePutSubj, data)
	if err != nil {
		return err
	}

	var resp genrepb.CreateResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (gi GenreAPI) Update(ctx *fiber.Ctx) error {
	var req genrepb.UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := gi.conn.RequestWithContext(ctx.Context(), genreUpdateSubj, data)
	if err != nil {
		return err
	}

	var resp genrepb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}

func (gi GenreAPI) Delete(ctx *fiber.Ctx) error {
	var req genrepb.DeleteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	respMsg, err := gi.conn.RequestWithContext(ctx.Context(), genreDeleteSubj, data)
	if err != nil {
		return err
	}

	var resp genrepb.EmptyResponse
	if err = proto.Unmarshal(respMsg.Data, &resp); err != nil {
		return err
	}

	return ctx.Status(int(resp.StatusCode)).JSON(&resp)
}
