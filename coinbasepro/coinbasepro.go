package coinbasepro

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type CoinbasePro struct {
	CsvFillsTXs   []CsvFillsTX
	CsvAccountTXs []CsvAccountTX
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *CoinbasePro {
	cbp := &CoinbasePro{}
	cbp.TXsByCategory = make(map[string]wallet.TXs)
	cbp.Sources = make(source.Sources)
	return cbp
}
