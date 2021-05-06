package revolut

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Revolut struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Revolut {
	revo := &Revolut{}
	revo.TXsByCategory = make(map[string]wallet.TXs)
	revo.Sources = make(source.Sources)
	return revo
}
