package kraken

// Kraken use specific asset names
func ReplaceAssets(assetToReplace string) string {
	switch assetToReplace {
	case "DOT.S":
		return "DOT"
	case "XETC":
		return "ETC"
	case "XETH":
		return "ETH"
	case "XLTC":
		return "LTC"
	case "XREP":
		return "REP"
	case "XXBT":
		return "BTC"
	case "XXDG":
		return "DOGE"
	case "XXRP":
		return "XRP"
	case "ZEUR":
		return "EUR"	
	default:
		return assetToReplace
	}
}