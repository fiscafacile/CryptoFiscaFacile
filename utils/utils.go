package utils

import (
	"strings"
)

func RemoveSymbol(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "@", "AROBASE"), ".", "POINT"), "-", "TIRET")
}

func AppendUniq(strs []string, str string) []string {
	found := false
	for _, s := range strs {
		if str == s {
			found = true
		}
	}
	if !found {
		return append(strs, str)
	}
	return strs
}
