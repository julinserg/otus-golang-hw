package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func compareFiles(file1, file2 string) bool {
	f1, err1 := ioutil.ReadFile(file1)

	if err1 != nil {
		fmt.Println(err1)
		return false
	}

	f2, err2 := ioutil.ReadFile(file2)

	if err2 != nil {
		fmt.Println(err2)
		return false
	}

	return bytes.Equal(f1, f2)
}

func TestCopy(t *testing.T) {
	var from, to string
	var limit, offset int64

	from = "testdata/input.txt"
	to = "out.txt"

	t.Run("full copy", func(t *testing.T) {
		defer os.Remove(to)
		err := Copy(from, to, offset, limit)
		require.Nil(t, err)
		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), srcFileStat.Size())
		require.True(t, compareFiles(to, "out.txt"))
	})

	t.Run("limit copy", func(t *testing.T) {
		defer os.Remove(to)
		limit = 10
		err := Copy(from, to, offset, limit)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), limit)
		require.True(t, compareFiles(to, "testdata/out_offset0_limit10.txt"))
	})

	t.Run("limit copy with limit over size", func(t *testing.T) {
		defer os.Remove(to)
		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		limit = srcFileStat.Size() + 10
		err = Copy(from, to, offset, limit)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), srcFileStat.Size())
		require.True(t, compareFiles(to, "out.txt"))
	})

	t.Run("offset and limit copy", func(t *testing.T) {
		defer os.Remove(to)
		limit = 1000
		offset = 100
		err := Copy(from, to, offset, limit)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), limit)
		require.True(t, compareFiles(to, "testdata/out_offset100_limit1000.txt"))
	})

	t.Run("offset and limit copy with limit over size", func(t *testing.T) {
		defer os.Remove(to)
		limit = 1000
		offset = 6000
		err := Copy(from, to, offset, limit)
		require.Nil(t, err)
		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), srcFileStat.Size()-offset)
		require.True(t, compareFiles(to, "testdata/out_offset6000_limit1000.txt"))
	})

	t.Run("error - offset copy with offset over size", func(t *testing.T) {
		defer os.Remove(to)
		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		limit = 0
		offset = srcFileStat.Size() + 1
		err = Copy(from, to, offset, limit)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("error - offset copy with offset equal size", func(t *testing.T) {
		defer os.Remove(to)
		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		limit = 0
		offset = srcFileStat.Size()
		err = Copy(from, to, offset, limit)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("error - copy not valid file", func(t *testing.T) {
		defer os.Remove(to)
		err := Copy("/dev/urandom", to, offset, limit)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("error - copy to root", func(t *testing.T) {
		defer os.Remove(to)
		err := Copy(from, "root/out.txt", offset, limit)
		require.NotNil(t, err)
	})

	t.Run("error - copy not exist file", func(t *testing.T) {
		defer os.Remove(to)
		err := Copy("/hgyt450/in.txt", to, offset, limit)
		require.NotNil(t, err)
	})

	t.Run("one byte copy - 1", func(t *testing.T) {
		defer os.Remove(to)

		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		offset = srcFileStat.Size() - 1
		limit = 0
		err = Copy(from, to, offset, limit)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), int64(1))
	})

	t.Run("one byte copy - 2", func(t *testing.T) {
		defer os.Remove(to)

		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		offset = srcFileStat.Size() - 1
		limit = 1
		err = Copy(from, to, offset, limit)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), int64(1))
	})

	t.Run("one byte copy - 3", func(t *testing.T) {
		defer os.Remove(to)

		srcFileStat, err := os.Stat(from)
		require.Nil(t, err)
		offset = srcFileStat.Size() - 1
		limit = 10
		err = Copy(from, to, offset, limit)
		require.Nil(t, err)
		dstFileStat, err := os.Stat(to)
		require.Nil(t, err)
		require.Equal(t, dstFileStat.Size(), int64(1))
	})
}
