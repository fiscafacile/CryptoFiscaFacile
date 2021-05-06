package binance

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Binance struct {
	csvTXs        []csvTX
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Binance {
	b := &Binance{}
	b.TXsByCategory = make(wallet.TXsByCategory)
	b.Sources = make(source.Sources)
	return b
}
