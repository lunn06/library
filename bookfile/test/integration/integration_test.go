//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	natscontainer "github.com/testcontainers/testcontainers-go/modules/nats"
	"testing"

	"github.com/lunn06/library/bookfile/client"
)

var boofileClient *client.Client

func TestMain(m *testing.M) {
	natsC := natsContainer()
	natsEndpoint, err := natsC.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	var cfg client.Config
	cfg.URL = fmt.Sprintf("nats://%s", natsEndpoint)

	boofileClient, err = client.New(cfg)
	if err != nil {
		panic(err)
	}

	m.Run()

	err = errors.Join(
		boofileClient.Close(),
		natsC.Terminate(context.Background()),
	)
	if err != nil {
		panic(err)
	}
}

func natsContainer() testcontainers.Container {
	ctx := context.Background()
	natsC, err := natscontainer.Run(ctx, "nats:2.11-alpine3.21")
	if err != nil {
		panic(err)
	}

	return natsC
}
