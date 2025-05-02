package proto

//go:generate protoc -I=../../../../api/book/ --go_out=. ../../../../api/book/author.proto ../../../../api/book/book.proto ../../../../api/book/genre.proto
