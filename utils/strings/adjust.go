package strings

import (
	"fmt"
	"strings"
)

func StringLeftJust(text string, filler string, size int) string {
	repeatSize := size - len(text)
	return fmt.Sprintf("%s%s", text, strings.Repeat(filler, repeatSize))
}

func StringRemovePrefix(text string, prefix string, trimSpaces bool) string {
	newString := ""

	// remove prefix
	if len(text) >= len(prefix) {
		newString = text[len(prefix):]
	}

	// trim space
	if trimSpaces && len(newString) > 1 {
		newString = strings.TrimSpace(newString)
	}

	return newString
}
