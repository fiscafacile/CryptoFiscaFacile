package uphold

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Uphold struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Uphold {
	uh := &Uphold{}
	uh.TXsByCategory = make(map[string]wallet.TXs)
	uh.Sources = make(source.Sources)
	return uh
}
