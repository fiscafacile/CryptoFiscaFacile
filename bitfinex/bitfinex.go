package bitfinex

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bitfinex struct {
	CsvTXs   []CsvTX
	Accounts wallet.Accounts
}

func New() *Bitfinex {
	bf := &Bitfinex{}
	bf.Accounts = make(map[string]wallet.TXs)
	return bf
}
