package fonts

import (
	"fmt"
	"unicode"
)

func GetCharacter(char rune, size string) ([]uint16, error) {
	if size == "1" {
		charData, ok := characters5x8[char]
		if !ok {
			return nil, fmt.Errorf("character '%c' in size '1' not found", char)
		}

		return charData, nil
	}

	if size == "3" {
		// Until I write all the lower letters in this bigger font
		if unicode.IsLower(char) && unicode.IsLetter(char) {
			char = unicode.ToUpper(char)
		}

		charData, ok := characters14x9[char]
		if !ok {
			return nil, fmt.Errorf("character '%c' in size '3' not found", char)
		}
		return charData, nil
	}

	if size == "2" {
		charData, ok := characters14x6[char]
		if !ok {
			return nil, fmt.Errorf("character '%c' in size '2' not found", char)
		}
		return charData, nil
	}

	return nil, fmt.Errorf("size '%s' not supported, must be '1', '2' or '3'", size)
}
