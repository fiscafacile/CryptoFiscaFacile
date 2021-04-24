package revolut

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Revolut struct {
	CsvTXs   []CsvTX
	Accounts wallet.Accounts
}

func New() *Revolut {
	revo := &Revolut{}
	revo.Accounts = make(map[string]wallet.TXs)
	return revo
}
