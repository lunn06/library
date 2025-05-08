package js

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	BookBucket     = "bookbucket"
	defaultTimeout = time.Minute * 5
)

func NewJetStream(url string) (jetstream.JetStream, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return jetstream.New(nc, jetstream.WithDefaultTimeout(defaultTimeout))
}

func NewObjectStore(js jetstream.JetStream) (jetstream.ObjectStore, error) {
	return js.CreateOrUpdateObjectStore(context.Background(), jetstream.ObjectStoreConfig{
		Bucket:      BookBucket,
		Storage:     jetstream.FileStorage,
		Compression: true,
	})
}
