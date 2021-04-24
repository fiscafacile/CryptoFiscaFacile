package bitfinex

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bitfinex struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *Bitfinex {
	bf := &Bitfinex{}
	bf.TXsByCategory = make(map[string]wallet.TXs)
	return bf
}
