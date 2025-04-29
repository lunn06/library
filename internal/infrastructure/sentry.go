package infrastructure

import "github.com/getsentry/sentry-go"

// SetupSentry init Sentry integration with pre-defined parameters.
func SetupSentry() {
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		SampleRate:       0.5,
		EnableTracing:    true,
	})
	if err != nil {
		panic(err)
	}
}
