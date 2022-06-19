package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const slashSymbol = `\`

func buildOneSymbol(prevSymbol rune, prevSymbolExist bool, strBuilder *strings.Builder) {
	if prevSymbolExist {
		strBuilder.WriteRune(prevSymbol)
	}
}

func buildSequenceSymbols(symbol rune, prevSymbol rune, prevSymbolExist bool, strBuilder *strings.Builder) error {
	if !prevSymbolExist {
		return ErrInvalidString
	}

	digit, err := strconv.Atoi(string(symbol))
	if err != nil {
		return err
	}

	if digit != 0 {
		for i := 0; i < digit; i++ {
			strBuilder.WriteRune(prevSymbol)
		}
	}
	return nil
}

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	var resultStrBuilder strings.Builder
	var prevSymbol rune
	prevSymbolExist := false
	slashSymbolExist := false
	for _, symbol := range str {
		if !unicode.IsDigit(symbol) && string(prevSymbol) == slashSymbol &&
			string(symbol) != slashSymbol && !slashSymbolExist {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(symbol) {
			if string(prevSymbol) == slashSymbol && !slashSymbolExist {
				prevSymbol = symbol
				continue
			}
			err := buildSequenceSymbols(symbol, prevSymbol, prevSymbolExist, &resultStrBuilder)
			if err != nil {
				return "", err
			}
			prevSymbolExist = false
		} else {
			if string(symbol) == slashSymbol && string(prevSymbol) == slashSymbol && !slashSymbolExist {
				slashSymbolExist = true
				continue
			}
			buildOneSymbol(prevSymbol, prevSymbolExist, &resultStrBuilder)
			prevSymbol = symbol
			prevSymbolExist = true
			slashSymbolExist = false
		}
	}
	if prevSymbolExist {
		resultStrBuilder.WriteRune(prevSymbol)
	}

	return resultStrBuilder.String(), nil
}
