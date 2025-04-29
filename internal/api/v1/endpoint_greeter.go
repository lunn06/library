package v1

import (
	"net/http"

	"github.com/exepirit/go-template/internal/service/greeter"
	"github.com/gin-gonic/gin"
)

func NewGreeterEndpoints(svc greeter.IService) *GreeterEndpoints {
	return &GreeterEndpoints{
		svc: svc,
	}
}

type GreeterEndpoints struct {
	svc greeter.IService
}

func (e GreeterEndpoints) Bind(router gin.IRouter) {
	router.POST("/greet", e.Greet)
}

func (e GreeterEndpoints) Greet(ctx *gin.Context) {
	var request greeter.GreetingRequest
	if err := ctx.BindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	greeting := e.svc.Greet(ctx, request)

	ctx.JSON(http.StatusOK, greeting)
}
