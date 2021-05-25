package monero

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Monero struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *Monero {
	xmr := &Monero{}
	xmr.TXsByCategory = make(map[string]wallet.TXs)
	return xmr
}
