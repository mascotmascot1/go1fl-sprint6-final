package service

import (
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// IsMorse checks whether the given string contains only valid Morse code symbols.
//
// Valid symbols are '.', '-', space, and '/'.
func IsMorse(phrase string) bool {
	for _, symbol := range phrase {
		if symbol != '.' && symbol != '-' && symbol != ' ' && symbol != '/' {
			return false
		}
	}
	return true
}

// Convert detects whether the input string is Morse code or plain text and converts it accordingly.
//
// If the input appears to be Morse code, it will be decoded to text.
// Otherwise, it will be encoded into Morse code.
// The function trims whitespace and returns an error if the input is empty.
func Convert(phrase string) (string, error) {
	phrase = strings.TrimSpace(phrase)
	if phrase == "" {
		return "", fmt.Errorf("input string is empty")
	}
	if IsMorse(phrase) {
		return morse.ToText(phrase), nil
	}
	return morse.ToMorse(phrase), nil
}
