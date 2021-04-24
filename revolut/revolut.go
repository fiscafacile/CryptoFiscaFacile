package revolut

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Revolut struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *Revolut {
	revo := &Revolut{}
	revo.TXsByCategory = make(map[string]wallet.TXs)
	return revo
}
