package domain

type Author struct {
	ID          int
	Name        string
	Description string
	Books       []Book
}
