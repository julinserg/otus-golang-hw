package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func createAndFillEnvFile(dpath, key string, value []byte) error {
	fname := filepath.Join(dpath, key)
	err := os.WriteFile(fname, value, 0o660)
	return err
}

func TestReadDir(t *testing.T) {
	dname, err := os.MkdirTemp("", "env_reader_test")
	defer os.RemoveAll(dname)
	require.Nil(t, err)

	err = createAndFillEnvFile(dname, "ENVKEY1", []byte("12"))
	require.Nil(t, err)
	err = createAndFillEnvFile(dname, "ENVKEY2", []byte("abcd     "))
	require.Nil(t, err)
	err = createAndFillEnvFile(dname, "ENVKEY3", []byte("abcd	"))
	require.Nil(t, err)
	err = createAndFillEnvFile(dname, "ENVKEY4", []byte("abcd\x00dcba"))
	require.Nil(t, err)
	err = createAndFillEnvFile(dname, "ENV=KEY5", []byte("12"))
	require.Nil(t, err)
	err = createAndFillEnvFile(dname, "ENVKEY6", []byte{})
	require.Nil(t, err)
	err = createAndFillEnvFile(dname, "ENVKEY7", []byte("       \ndcba"))
	require.Nil(t, err)

	env, err := ReadDir(dname)
	require.Nil(t, err)
	require.Equal(t, env["ENVKEY1"].Value, "12")
	require.Equal(t, env["ENVKEY2"].Value, "abcd")
	require.Equal(t, env["ENVKEY3"].Value, "abcd")
	require.Equal(t, env["ENVKEY4"].Value, "abcd\ndcba")
	require.Equal(t, env["ENV=KEY5"].Value, "")
	require.Equal(t, env["ENVKEY6"].Value, "")
	require.Equal(t, env["ENVKEY6"].NeedRemove, true)
	require.Equal(t, env["ENVKEY7"].Value, "")
	require.Equal(t, env["ENVKEY7"].NeedRemove, false)

	_, err = ReadDir("/root")
	require.NotNil(t, err)

	dnameEmpty, err := os.MkdirTemp("", "env_reader_test")
	defer os.RemoveAll(dnameEmpty)
	require.Nil(t, err)
	env, err = ReadDir(dnameEmpty)
	require.Nil(t, err)
	require.Equal(t, len(env), 0)
}
