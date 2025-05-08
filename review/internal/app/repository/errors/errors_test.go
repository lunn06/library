package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsErrNotFound(t *testing.T) {
	assert.True(t, IsErrNotFound(ErrNotFound{Inner: errors.New("outer: not found")}))
	assert.False(t, IsErrNotFound(errors.New("not not found")))
	assert.True(t, IsErrNotFound(
		fmt.Errorf("wrapped: %w", ErrNotFound{Inner: errors.New("outer: not found")}),
	))
}
