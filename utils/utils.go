package utils

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
