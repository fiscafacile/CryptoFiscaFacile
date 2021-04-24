package binance

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Binance struct {
	csvTXs         []csvTX
	csvExtendedTXs []csvExtendedTX
	Accounts       wallet.Accounts
}

func New() *Binance {
	b := &Binance{}
	b.Accounts = make(map[string]wallet.TXs)
	return b
}
