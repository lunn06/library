//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	natscontainer "github.com/testcontainers/testcontainers-go/modules/nats"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/protobuf/proto"

	natsapi "github.com/lunn06/library/review/internal/api/nats"
	"github.com/lunn06/library/review/internal/app/repository"
	"github.com/lunn06/library/review/internal/app/service"
	"github.com/lunn06/library/review/internal/config"
	"github.com/lunn06/library/review/internal/infrastructure/db/postgres"
)

const reqTimeout = time.Second

const (
	reviewGetAllByBookIDSubj = "gateway.getAllByBookID"
	reviewGetSubj            = "gateway.get"
	reviewPutSubj            = "gateway.put"
	reviewUpdateSubj         = "gateway.update"
	reviewDeleteSubj         = "gateway.delete"
)

const (
	postgresUser     = "test-gateway-user"
	postgresPassword = "test-gateway-password"
	postgresDb       = "test-gateway-db"
	postgresSslMode  = "disable"
)

var nc *nats.Conn

func TestMain(m *testing.M) {
	natsC := natsContainer()
	natsEndpoint, err := natsC.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	postgresC := postgresContainer()
	postgresEndpoint, err := postgresC.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	cfg.Nats = natsapi.Config{
		URL: fmt.Sprintf("nats://%s", natsEndpoint),
	}
	cfg.Postgres.URL = postgresEndpoint
	cfg.Postgres.User = postgresUser
	cfg.Postgres.Password = postgresPassword
	cfg.Postgres.DB = postgresDb
	cfg.Postgres.SslMode = postgresSslMode

	entClient, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	reviewRepo := repository.NewEntReviewRepo(entClient)
	reviewService := service.NewReviewService(reviewRepo)
	reviewConsumer := natsapi.NewReviewConsumer(reviewService)

	nc, err = nats.Connect(cfg.Nats.URL)

	if err = natsapi.RegisterReviewConsumer(nc, reviewConsumer); err != nil {
		panic(err)
	}

	m.Run()

	err = errors.Join(
		nc.Drain(),
		natsC.Terminate(context.Background()),
		postgresC.Terminate(context.Background()),
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

func postgresContainer() testcontainers.Container {
	ctx := context.Background()
	postgresC, err := postgrescontainer.Run(ctx,
		"postgres:17-alpine3.21",
		postgrescontainer.WithUsername(postgresUser),
		postgrescontainer.WithPassword(postgresPassword),
		postgrescontainer.WithDatabase(postgresDb),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		panic(err)
	}

	return postgresC
}

func request(subj string, req, resp proto.Message) error {
	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	resMsg, err := nc.Request(subj, data, reqTimeout)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(resMsg.Data, resp)
	if err != nil {
		return err
	}

	return nil
}
