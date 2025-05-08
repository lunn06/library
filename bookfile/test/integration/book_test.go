//go:build integration

package integration

import (
	"bytes"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	servicerrors "github.com/lunn06/library/bookfile/internal/app/service/errors"
	"github.com/lunn06/library/bookfile/internal/domain"

	_ "embed"
)

const testFileName = "test1.txt"

//go:embed test.dat
var testFile []byte

func BenchmarkBookFileCreate(b *testing.B) {
	for b.Loop() {
		_, err := boofileClient.Create(b.Context(), testFileName, bytes.NewBuffer(testFile))
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkBookFileGetMultipleFile(b *testing.B) {
	uuids := make([]uuid.UUID, b.N)
	for i := 0; i < b.N; i++ {
		bookUUID, err := boofileClient.Create(b.Context(), testFileName, bytes.NewBuffer(testFile))
		if err != nil {
			b.Fatal(err)
		}

		uuids[i] = bookUUID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := boofileClient.Get(b.Context(), uuids[i])
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkBookFileGetOneFile(b *testing.B) {
	bookUUID, err := boofileClient.Create(b.Context(), testFileName, bytes.NewBuffer(testFile))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = boofileClient.Get(b.Context(), bookUUID)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func TestBookFileCreate(t *testing.T) {
	// Put bookfile
	bookUUID, err := boofileClient.Create(t.Context(), testFileName, bytes.NewBuffer(testFile))
	require.NoError(t, err)
	//////////////

	// Check putted bookfile
	book, err := boofileClient.Get(t.Context(), bookUUID)
	require.NoError(t, err)

	assert.Equal(t, bookUUID, book.UUID)
	assert.Equal(t, testFileName, book.FileName)
	assert.Equal(t, testFile, book.Buffer.Bytes())
}

func TestBookFileDelete(t *testing.T) {
	// Put bookfile
	bookUUID, err := boofileClient.Create(t.Context(), testFileName, bytes.NewBuffer(testFile))
	require.NoError(t, err)
	//////////////

	// Check putted bookfile
	book, err := boofileClient.Get(t.Context(), bookUUID)
	require.NoError(t, err)

	assert.Equal(t, testFileName, book.FileName)
	assert.Equal(t, testFile, book.Buffer.Bytes())
	//////////////

	// Delete bookfile
	err = boofileClient.Delete(t.Context(), bookUUID)
	require.NoError(t, err)
	//////////////

	// Check deleted bookfile not found
	book, err = boofileClient.Get(t.Context(), bookUUID)

	assert.Equal(t, domain.Book{}, book)
	assert.Error(t, servicerrors.ErrResourceNotFound{}, err)
}
