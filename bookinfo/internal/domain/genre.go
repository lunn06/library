package domain

type Genre struct {
	ID          int
	Title       string
	Description string
	Books       []Book
}
