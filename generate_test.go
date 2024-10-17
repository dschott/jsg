package main

import (
	"go/format"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Generator(t *testing.T) {
	var generator Generator

	file := File{
		Version: "0.0.0",
		Pkg:     "fooPackage",
		Types: []*Type{
			{
				Name: "FooType",
			}, {
				Name: "BarType",
			},
		},
	}

	var sb strings.Builder
	err := generator.Generate(&sb, &file)
	require.NoError(t, err)

	got := sb.String()

	formatted, err := format.Source([]byte(got))
	require.NoError(t, err, "Got:\n"+got)

	require.Equal(t, string(formatted), got)
}
