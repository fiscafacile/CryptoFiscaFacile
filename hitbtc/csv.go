package hitbtc

import (
	"strings"
)

func csvCurrencyCure(c string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(c, "BCHSV", "BSV"), "BCHABC", "BCH"), "BCCF", "BCHOLD"))
}
