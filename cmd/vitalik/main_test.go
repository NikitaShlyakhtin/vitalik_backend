package main

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(getFxOptions())
	require.NoError(t, err)
}
