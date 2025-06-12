package domain

import (
	"io"

	"github.com/google/uuid"
)

type Book struct {
	UUID           uuid.UUID
	FileName       string
	FileReadCloser io.ReadCloser
}

func (b *Book) Read(p []byte) (int, error) {
	return b.FileReadCloser.Read(p)
}

func (b *Book) Close() error {
	return b.FileReadCloser.Close()
}
