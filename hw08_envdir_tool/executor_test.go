package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{"ENV1": EnvValue{"val1", false}, "ENV2": EnvValue{"", true}, "ENV3": EnvValue{"", false}}
	retCode := RunCmd([]string{"./testdata/echo.sh"}, env)
	require.Equal(t, retCode, 0)

	retCode = RunCmd([]string{"/root/admin.sh"}, env)
	require.Equal(t, retCode, 1)

	retCode = RunCmd([]string{""}, env)
	require.Equal(t, retCode, 1)

	retCode = RunCmd(nil, env)
	require.Equal(t, retCode, 1)

	env = Environment{"ENV1": EnvValue{"abcd\x00dcba", false}}
	retCode = RunCmd([]string{"./testdata/echo.sh"}, env)
	require.Equal(t, retCode, 1)

	// ...............................

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	env = Environment{
		"HELLO": EnvValue{"\"hello\"", false},
		"BAR":   EnvValue{"bar", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"UNSET": EnvValue{"", true},
		"EMPTY": EnvValue{"", false},
	}

	os.Setenv("ADDED", "from original env")
	retCode = RunCmd([]string{"./testdata/echo.sh", "arg1=1", "arg2=2"}, env)
	require.Equal(t, retCode, 0)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	s := string(out)
	expected := "HELLO is (\"hello\")\n" +
		"BAR is (bar)\n" +
		"FOO is (   foo\nwith new line)\n" +
		"UNSET is ()\n" +
		"ADDED is (from original env)\n" +
		"EMPTY is ()\n" +
		"arguments are arg1=1 arg2=2\n"
	require.Equal(t, s, expected)
}
