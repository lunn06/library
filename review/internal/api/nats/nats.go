package nats

import (
	"log/slog"
	"net/http"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	reviewpb "github.com/lunn06/library/review/internal/api/proto/review"
)

func NewConnection(cfg Config) (*nats.Conn, error) {
	return nats.Connect(cfg.URL)
}

func NewResponder(msg *nats.Msg) Responder {
	return Responder{msg: msg}
}

type Responder struct {
	msg *nats.Msg
}

func (r Responder) Respond(req proto.Message) {
	data, err := proto.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal request", "err", err)
		data, err = proto.Marshal(&reviewpb.EmptyResponse{
			StatusCode: http.StatusInternalServerError,
		})
		if err != nil {
			slog.Error("Failed to marshal request", "err", err)
			_ = r.msg.Nak()
		}

		_ = r.msg.Respond(data)
	}

	_ = r.msg.Respond(data)
}
