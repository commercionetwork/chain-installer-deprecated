package utils

import (
	"fmt"
	"os"
	"strings"
)

// CheckError correctly displays the given error to the user, and exits the program with non-zero code
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// GetUserHome allow to return the string representing the path to the user home directory
func GetUserHome() string {
	home, err := os.UserHomeDir()
	CheckError(err)
	return home
}

// ReplaceLast returns the `original` string having the last occurrence of the specified `old` string replaced with the
// specified `replace` string
func ReplaceLast(original, old, replace string) string {
	i := strings.LastIndex(original, old)

	if i >= 0 {
		return original[:i] + strings.Replace(original[i:], old, replace, 1)
	} else {
		return original
	}
}
