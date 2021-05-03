package kraken

// Kraken use specific asset names
func ReplaceAssets(assetToReplace string) string {
	switch assetToReplace {
	case "DOT.S":
		return "DOT"
	case "XBT":
		return "BTC"
	case "XDG":
		return "DOGE"
	default:
		return assetToReplace
	}
}
