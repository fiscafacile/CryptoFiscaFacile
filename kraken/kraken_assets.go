package kraken

import "strings"

// Kraken use specific asset names
func ReplaceAssets(assetToReplace string) string {
	assetRplcr := strings.NewReplacer(
		"ADA.S", "ADA",
		"ATOM.S", "ATOM",
		"DOT.S", "DOT",
		"ETH2", "ETH",
		"ETH2.S", "ETH",
		"EUR.HOLD", "EUR",
		"EUR.M", "EUR",
		"FLOW.S", "FLOW",
		"FLOWH", "FLOW",
		"FLOWH.S", "FLOW",
		"KAVA.S", "KAVA",
		"KFEE", "FEE",
		"KSM.S", "KSM",
		"USD.HOLD", "USD",
		"USD.M", "USD",
		"XBT", "BTC",
		"XBT.M", "BTC",
		"XETC", "ETC",
		"XETH", "ETH",
		"XLTC", "LTC",
		"XMLN", "MLN",
		"XREP", "REP",
		"XTZ", "XTZ",
		"XTZ.S", "XTZ",
		"XXBT", "BTC",
		"XXDG", "DOGE",
		"XXLM", "XLM",
		"XXMR", "XMR",
		"XXRP", "XRP",
		"XZEC", "ZEC",
		"ZAUD", "AUD",
		"ZCAD", "CAD",
		"ZEUR", "EUR",
		"ZGBP", "GBP",
		"ZJPY", "JPY",
		"ZRX", "ZRX",
		"ZUSD", "USD",
	)
	return assetRplcr.Replace(assetToReplace)
}
