package fonts

import (
	"fmt"
	"unicode"
)

func GetCharacter(char rune, size string) ([]uint16, error) {
	if size == "small" {
		charData, ok := characters5x8[char]
		if !ok {
			return nil, fmt.Errorf("character '%c' in size 'small' not found", char)
		}

		return charData, nil
	}

	if size == "large" {
		// Until I write all the lower letters in this bigger font
		if unicode.IsLower(char) && unicode.IsLetter(char) {
			char = unicode.ToUpper(char)
		}

		charData, ok := characters14x9[char]
		if !ok {
			return nil, fmt.Errorf("character '%c' in size 'large' not found", char)
		}
		return charData, nil
	}

	return nil, fmt.Errorf("size '%s' not supported, must be 'small' or 'large'", size)
}
