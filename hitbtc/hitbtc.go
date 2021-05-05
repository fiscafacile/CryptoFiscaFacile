package hitbtc

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type HitBTC struct {
	csvTradeTXs       []csvTradeTX
	csvTransactionTXs []csvTransactionTX
	TXsByCategory     wallet.TXsByCategory
}

func New() *HitBTC {
	hb := &HitBTC{}
	hb.TXsByCategory = make(wallet.TXsByCategory)
	return hb
}
