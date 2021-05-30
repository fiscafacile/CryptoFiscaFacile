package coinbase

import (
	"strings"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Coinbase struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Coinbase {
	cb := &Coinbase{}
	cb.TXsByCategory = make(map[string]wallet.TXs)
	cb.Sources = make(source.Sources)
	return cb
}

func ReplaceAssets(assetToReplace string) string {
	assetRplcr := strings.NewReplacer(
		"CGLD", "CELO",
	)
	return assetRplcr.Replace(assetToReplace)
}
