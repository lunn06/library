package greeter

import (
	"context"
	"fmt"
)

type GreetingRequest struct {
	Name string `form:"name" binding:"required"`
}

type Greeting struct {
	Text string `json:"text"`
}

type IService interface {
	Greet(ctx context.Context, request GreetingRequest) Greeting
}

func NewService() IService {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (serviceImpl) Greet(ctx context.Context, request GreetingRequest) Greeting {
	text := fmt.Sprintf("Hello, %s!", request.Name)
	return Greeting{
		Text: text,
	}
}
