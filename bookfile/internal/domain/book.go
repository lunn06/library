package domain

import (
	"bytes"

	"github.com/google/uuid"
)

type Book struct {
	UUID     uuid.UUID
	FileName string
	Buffer   *bytes.Buffer
}
