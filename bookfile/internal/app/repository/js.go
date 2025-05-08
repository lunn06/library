package repository

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"

	repoerrors "github.com/lunn06/library/bookfile/internal/app/repository/errors"
	"github.com/lunn06/library/bookfile/internal/domain"
)

const fileNameKey = "fileName"

var _ BookRepo = (*JsBookRepo)(nil)

func NewJsBookRepo(store jetstream.ObjectStore) *JsBookRepo {
	return &JsBookRepo{
		ObjectStore: store,
	}
}

type JsBookRepo struct {
	jetstream.ObjectStore
}

func (jbr *JsBookRepo) Get(ctx context.Context, bookUUID uuid.UUID) (domain.Book, error) {
	obj, err := jbr.ObjectStore.Get(ctx, bookUUID.String())
	if err != nil {
		if errors.Is(err, jetstream.ErrObjectNotFound) {
			return domain.Book{}, repoerrors.ErrNotFound{Inner: err}
		}
		return domain.Book{}, err
	}
	defer obj.Close()

	objInfo, err := obj.Info()
	if err != nil {
		return domain.Book{}, err
	}

	bytesBuf := make([]byte, objInfo.Size)
	if _, err = obj.Read(bytesBuf); err != nil {
		return domain.Book{}, err
	}

	fileName := objInfo.Metadata[fileNameKey]

	return domain.Book{
		UUID:     bookUUID,
		FileName: fileName,
		Buffer:   bytes.NewBuffer(bytesBuf),
	}, nil
}

func (jbr *JsBookRepo) Put(ctx context.Context, book domain.Book) error {
	meta := jetstream.ObjectMeta{
		Name: book.UUID.String(),
		Metadata: map[string]string{
			fileNameKey: book.FileName,
		},
	}

	_, err := jbr.ObjectStore.Put(ctx, meta, book.Buffer)
	return err
}

func (jbr *JsBookRepo) Delete(ctx context.Context, bookUUID uuid.UUID) error {
	return jbr.ObjectStore.Delete(ctx, bookUUID.String())
}
