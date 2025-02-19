package main

import (
	"testing"

	_ "embed"

	"github.com/stretchr/testify/require"
)

//go:embed source_test.txt
var ref string

func TestTags(t *testing.T) {
	err := Parse(ref, "newTag")
	require.NoError(t, err)
}

func TestGenerateTag(t *testing.T) {
	newTag := generateTag("TotoTata", "c", `a:"a" b:"b"`)
	require.Equal(t, "`a:\"a\" b:\"b\" c:\"toto_tata\"`", newTag)
}
