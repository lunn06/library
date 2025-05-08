package middleware

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

type Middleware func(handler nats.MsgHandler) nats.MsgHandler

func With(base nats.MsgHandler, mws ...Middleware) nats.MsgHandler {
	for _, mw := range mws {
		base = mw(base)
	}

	return base
}

func Recover() Middleware {
	return func(handler nats.MsgHandler) nats.MsgHandler {
		return func(msg *nats.Msg) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered",
						slog.Any("panic", r),
					)
					_ = msg.Nak()
				}
			}()
			handler(msg)
		}
	}
}

func Logger(logger *slog.Logger) Middleware {
	return func(handler nats.MsgHandler) nats.MsgHandler {
		return func(msg *nats.Msg) {
			start := time.Now()
			handler(msg)
			handleTime := time.Now().Sub(start)

			logger.Info("Received message",
				slog.Time("start", start),
				slog.String("subject", msg.Subject),
				slog.Duration("perf", handleTime),
				slog.Any("data", msg.Data),
			)
		}
	}
}
