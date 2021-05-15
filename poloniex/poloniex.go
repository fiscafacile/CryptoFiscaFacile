package poloniex

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Poloniex struct {
	csvDepositsTXs      []csvDepositsTX
	csvDistributionsTXs []csvDistributionsTX
	csvTradesTXs        []csvTradesTX
	csvWithdrawalsTXs   []csvWithdrawalsTX
	TXsByCategory       wallet.TXsByCategory
	Sources             source.Sources
}

func New() *Poloniex {
	pl := &Poloniex{}
	pl.TXsByCategory = make(wallet.TXsByCategory)
	pl.Sources = make(source.Sources)
	return pl
}
