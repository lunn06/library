package api

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lunn06/library/bookfile/client"
	"github.com/lunn06/library/gateway/internal/api/nats"

	bookfilepb "github.com/lunn06/library/gateway/internal/api/proto/bookfile"
)

//type DeleteResponse struct {
//	StatusCode int `json:"status_code"`
//}
//
//type DeleteRequest struct {
//	BookUUID string `json:"uuid"`
//}
//
//type GetRequest struct {
//	BookUUID string `params:"uuid"`
//}
//
//type GetResponse struct {
//	BookUUID string `json:"uuid"`
//	FileName string `json:"file_name"`
//	File     []byte `json:"file"`
//}
//
//type CreateRequest struct {
//	FileName string `json:"file_name"`
//	File     []byte `json:"file"`
//}
//
//type CreateResponse struct {
//	BookUUID   string `json:"uuid"`
//	StatusCode int    `json:"status_code"`
//}

func NewBookFileClient(config nats.Config) (*client.Client, error) {
	return client.New(client.Config{URL: config.URL})
}

func NewBookFileAPI(client *client.Client) BookFileAPI {
	return BookFileAPI{client: client}
}

type BookFileAPI struct {
	client *client.Client
}

func (bi BookFileAPI) Register(router fiber.Router) {
	router.
		Get("/book/file/:uuid", bi.Get).
		Post("/book/file", bi.Put).
		Delete("/book/file/:uuid", bi.Delete)
}

func (bi BookFileAPI) Get(ctx *fiber.Ctx) error {
	var params struct {
		UUID string `params:"uuid"`
	}
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	bookUUID, err := uuid.Parse(params.UUID)
	if err != nil {
		return err
	}

	book, err := bi.client.Get(ctx.Context(), bookUUID)
	if err != nil {
		return err
	}
	defer book.Close()

	//return ctx.SendStream(book.FileReader)

	bookBytes := &bytes.Buffer{}
	if _, err = bookBytes.ReadFrom(&book); err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(&bookfilepb.GetResponse{
		BookUuid: book.UUID.String(),
		FileName: book.FileName,
		File:     bookBytes.Bytes(),
	})
}

func (bi BookFileAPI) Put(ctx *fiber.Ctx) error {
	//var req CreateRequest
	//if err := ctx.BodyParser(&req); err != nil {
	//	return err
	//}

	mp, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	fileData, ok := mp.Value["book"]
	if !ok {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	fileName, ok := mp.Value["filename"]
	if !ok {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	bookUUID, err := bi.client.Create(ctx.Context(),
		fileName[0],
		ioutils.NewReadCloserWrapper(
			strings.NewReader(fileData[0]),
			func() error {
				return nil
			},
		),
	)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(&bookfilepb.CreateResponse{
		BookUuid: bookUUID.String(),
	})
}

func (bi BookFileAPI) Delete(ctx *fiber.Ctx) error {
	var params struct {
		UUID string `params:"uuid"`
	}
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	bookUUID, err := uuid.Parse(params.UUID)
	if err != nil {
		return err
	}

	if err = bi.client.Delete(ctx.Context(), bookUUID); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(&bookfilepb.DeleteResponse{
		StatusCode: http.StatusOK,
	})
}
