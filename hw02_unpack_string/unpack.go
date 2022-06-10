package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	var resultStrBuilder strings.Builder
	var prevS rune
	prevSExist := false
	for _, s := range str {
		if unicode.IsDigit(s) {
			if !prevSExist {
				return "", ErrInvalidString
			}

			digit, err := strconv.Atoi(string(s))
			if err != nil {
				return "", err
			}

			if digit != 0 {
				for i := 0; i < digit; i++ {
					resultStrBuilder.WriteRune(prevS)
				}
			}
			prevSExist = false
		} else {
			if prevSExist {
				resultStrBuilder.WriteRune(prevS)
			}
			prevS = s
			prevSExist = true
		}
	}
	if prevSExist {
		resultStrBuilder.WriteRune(prevS)
	}

	return resultStrBuilder.String(), nil
}
