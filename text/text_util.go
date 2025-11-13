package text

import (
	"strings"
)

const (
	CHAR_COLON         = ":"
	CHAR_TILDE         = "~"
	CHAR_FORWARD_SLASH = "/"
	CHAR_BACK_SLASH    = "\\"
	CHAR_FULL_STOP     = "."
	CHAR_ASTERIX       = "*"
	HOME_DIR_SHORTHAND = CHAR_TILDE + CHAR_FORWARD_SLASH
	EMPTY              = ""
	WHITE_SPACE        = " "
	COLON              = ":"
)

// StringBlank reports whether the string contains only whitespace characters.
// It treats an empty string as blank.
func StringBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// StringNotBlank is the negation of StringBlank.
func StringNotBlank(s string) bool {
	return !StringBlank(s)
}

// ListContains reports whether any element in the slice contains arg as a substring.
func ListContains(arguments []string, arg string) bool {
	for _, argItem := range arguments {
		if strings.Contains(argItem, arg) {
			return true
		}
	}
	return false
}

// GetArg returns the value for a flag of the form key=value.
// It scans the slice and returns the first matching value where the element
// starts with arg + "=". If no such element exists, it returns EMPTY.
func GetArg(arguments []string, arg string) string {
	prefix := arg + "="
	for _, argItem := range arguments {
		if strings.HasPrefix(argItem, prefix) {
			return strings.TrimPrefix(argItem, prefix)
		}
	}
	return EMPTY
}

// EqualsIgnoreCase compares two strings for equality, ignoring case.
func EqualsIgnoreCase(textArg string, anotherTextArg string) bool {
	return strings.EqualFold(textArg, anotherTextArg)
}

// Equals
// Checks string Equality including case sensitivity
func Equals(textArg string, anotherTextArg string) bool {
	return textArg == anotherTextArg
}

// NotEquals reports whether two strings are not equal (case sensitive).
func NotEquals(textArg string, anotherTextArg string) bool {
	return textArg != anotherTextArg
}

// NotEqualsIgnoreCase reports whether two strings are not equal, ignoring case.
func NotEqualsIgnoreCase(textArg string, anotherTextArg string) bool {
	return !EqualsIgnoreCase(textArg, anotherTextArg)
}

// Trim returns textValue with all leading and trailing white space removed.
func Trim(textValue string) string {
	return strings.TrimSpace(textValue)
}
