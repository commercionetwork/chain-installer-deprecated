package utils

import (
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetUserHome() string {
	home, err := os.UserHomeDir()
	CheckError(err)
	return home
}

func ReplaceLast(original, old, replace string) string {
	i := strings.LastIndex(original, old)

	if i >= 0 {
		return original[:i] + strings.Replace(original[i:], old, replace, 1)
	} else {
		return original
	}
}
