package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestFxApp(t *testing.T) {
	require.NoError(t, fx.ValidateApp(Module))
}
