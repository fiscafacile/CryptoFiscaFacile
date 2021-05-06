package bitfinex

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bitfinex struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Bitfinex {
	bf := &Bitfinex{}
	bf.TXsByCategory = make(wallet.TXsByCategory)
	bf.Sources = make(source.Sources)
	return bf
}
