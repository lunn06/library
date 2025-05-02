package domain

type Book struct {
	ID          int
	UserID      int
	Title       string
	Description string
	BookURL     string
	CoverURL    *string
	Authors     []Author
	Genres      []Genre
}
