package server

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/slog-fiber"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func NewServer(cfg Config) *Server {
	app := fiber.New(fiber.Config{
		Prefork:               cfg.Prefork,
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		IdleTimeout:           cfg.IdleTimeout,
		EnablePrintRoutes:     cfg.Debug,
		DisableStartupMessage: !cfg.Debug,
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return protojson.Marshal(v.(proto.Message))
		},
	})

	app.Use(slogfiber.New(slog.Default()))

	return &Server{
		app:  app,
		addr: cfg.Address,
	}
}

type Server struct {
	app  *fiber.App
	addr string
}

func (s Server) Router() fiber.Router {
	return s.app
}

func (s Server) Listen() error {
	return s.app.Listen(s.addr)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
