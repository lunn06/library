package postgres

type Config struct {
	URL      string `default:"localhost:5432"`
	User     string `required:"true"`
	Password string `required:"true"`
	DB       string `required:"true"`
	SslMode  string `default:"disable"`
}
