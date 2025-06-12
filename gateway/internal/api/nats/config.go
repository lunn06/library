package nats

type Config struct {
	URL string `default:"nats://127.0.0.1:4222"`
}
