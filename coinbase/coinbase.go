package coinbase

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Coinbase struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *Coinbase {
	cb := &Coinbase{}
	cb.TXsByCategory = make(map[string]wallet.TXs)
	return cb
}
