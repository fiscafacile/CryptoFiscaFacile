package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

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

func GetUniqueID(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}
