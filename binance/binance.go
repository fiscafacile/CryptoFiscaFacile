package binance

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Binance struct {
	csvTXs         []csvTX
	csvExtendedTXs []csvExtendedTX
	TXsByCategory  wallet.TXsByCategory
}

func New() *Binance {
	b := &Binance{}
	b.TXsByCategory = make(map[string]wallet.TXs)
	return b
}
