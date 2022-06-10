package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "aaa0b1", expected: "aab"},
		{input: "aaa0b2", expected: "aabb"},
		{input: "a0b1c2d3", expected: "bccddd"},
		{input: "ш0щ1я2ы3", expected: "щяяыыы"},
		{input: "j0s1я2ы3", expected: "sяяыыы"},
		{input: "¹º!0*1¹º×2¶3½4", expected: "¹º*¹º××¶¶¶½½½½"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwe\\\3\\`, expected: `qwe\3\`},
		{input: `qwe\\\\3`, expected: `qwe\\\\`},
		{input: `qwe\\\\`, expected: `qwe\\`},
		{input: `\\\\`, expected: `\\`},
		{input: `\\\\\\\\`, expected: `\\\\`},
		{input: `\\2`, expected: `\\`},
		{input: `\2aaa0b\3`, expected: `2aab3`},
		{input: `\2aaa0b\3\\\3`, expected: `2aab3\3`},
		{input: `\2aaa0b\3\\\3qwe\\\\ы3`, expected: `2aab3\3qwe\\ыыы`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
