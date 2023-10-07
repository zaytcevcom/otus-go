package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	var prevChar rune

	for i, char := range str {
		if unicode.IsDigit(char) && (unicode.IsDigit(prevChar) || i == 0) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(char) {
			s, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidString
			}

			result.WriteString(strings.Repeat(string(prevChar), s))
		} else if !unicode.IsDigit(char) && !unicode.IsDigit(prevChar) && i != 0 {
			result.WriteRune(prevChar)
		}

		if !unicode.IsDigit(char) && i == len(str)-1 {
			result.WriteRune(char)
		}

		prevChar = char
	}

	return result.String(), nil
}
