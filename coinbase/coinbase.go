package coinbase

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Coinbase struct {
	CsvTXs   []CsvTX
	Accounts wallet.Accounts
}

func New() *Coinbase {
	cb := &Coinbase{}
	cb.Accounts = make(map[string]wallet.TXs)
	return cb
}
