package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnakeCase(t *testing.T) {
	input := "ID"
	result := toSnakeCase(input)
	require.Equal(t, "id", result)
}
