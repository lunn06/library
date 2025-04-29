package ginmiddleware

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// SentryTracingMiddleware makes middleware, that record request as transaction.
func SentryTracingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transaction := sentry.StartTransaction(
			ctx, ctx.FullPath(),
			sentry.ContinueFromRequest(ctx.Request),
		)

		defer func(){
			// TODO: recover panic and finish transaction with error

			transaction.Status = matchTransactionStatus(ctx.Request.Response.StatusCode)
			transaction.Finish()
		}()

		ctx.Next()
	}
}

// SentryDefaultMiddleware is a wrapper around official Sentry library with little improvements.
func SentryDefaultMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic:         true,
	})
}

// matchTransactionStatus mathes HTTP status code to Sentry transaction status.
//
// More infofmation here: https://develop.sentry.dev/sdk/event-payloads/span/
func matchTransactionStatus(statusCode int) sentry.SpanStatus {
	if statusCode >= 200 && statusCode < 300 {
		return sentry.SpanStatusOK
	}

	switch statusCode {
	case 400:
		return sentry.SpanStatusFailedPrecondition
	case 401:
		return sentry.SpanStatusUnauthenticated
	case 403:
		return sentry.SpanStatusPermissionDenied
	case 404:
		return sentry.SpanStatusNotFound
	case 409:
		return sentry.SpanStatusAlreadyExists
	case 429:
		return sentry.SpanStatusResourceExhausted
	case 499:
		return sentry.SpanStatusCanceled
	case 500:
		return sentry.SpanStatusInternalError
	case 501:
		return sentry.SpanStatusUnimplemented
	case 503:
		return sentry.SpanStatusUnavailable
	case 504:
		return sentry.SpanStatusDeadlineExceeded
	default:
		return sentry.SpanStatusUnknown
	}
}
