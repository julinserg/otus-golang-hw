package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var from, to string
	var limit, offset int64

	from = "testdata/input.txt"
	to = "out.txt"
	err := Copy(from, to, offset, limit)
	require.Nil(t, err)

	srcFileStat, err := os.Stat(from)
	require.Nil(t, err)
	dstFileStat, err := os.Stat(to)
	require.Nil(t, err)
	require.Equal(t, srcFileStat.Size(), dstFileStat.Size())
}
