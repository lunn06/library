package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	bookpb "github.com/lunn06/library/book/internal/api/proto/book"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	req := bookpb.DeleteRequest{
		BookId: 5,
	}

	out, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}

	msg, err := nc.Request("book.delete", out, time.Second)
	if err != nil {
		panic(err)
	}

	var resp bookpb.EmptyResponse
	if err = proto.Unmarshal(msg.Data, &resp); err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}
