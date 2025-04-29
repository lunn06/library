package infrastructure

import (
	"net/http"

	"github.com/exepirit/go-template/internal/config"
	"github.com/exepirit/go-template/pkg/middleware/ginmiddleware"
	"github.com/exepirit/go-template/pkg/routing"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg config.Config) (Server, error) {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	ginHandler := gin.New()
	ginHandler.Use(
		gin.Recovery(),
		ginmiddleware.SentryDefaultMiddleware(),
		ginmiddleware.SentryTracingMiddleware(),
	)

	server := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: ginHandler,
	}

	return Server{
		Server: server,
	}, nil
}

type Server struct {
	*http.Server
}

func (srv Server) Bind(bindable routing.Bindable) {
	if bindable == nil {
		return
	}
	bindable.Bind(srv.Handler.(*gin.Engine))
}
