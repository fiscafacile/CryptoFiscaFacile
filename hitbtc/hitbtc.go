package hitbtc

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type HitBTC struct {
	api               api
	csvTradeTXs       []csvTradeTX
	csvTransactionTXs []csvTransactionTX
	done              chan error
	TXsByCategory     wallet.TXsByCategory
}

func New() *HitBTC {
	hb := &HitBTC{}
	hb.done = make(chan error)
	hb.TXsByCategory = make(wallet.TXsByCategory)
	return hb
}

func (hb *HitBTC) GetAPIAllTXs() {
	err := hb.api.getAllTXs()
	if err != nil {
		hb.done <- err
		return
	}
	hb.TXsByCategory.Add(hb.api.txsByCategory)
	hb.done <- nil
}

func (hb *HitBTC) WaitFinish() error {
	return <-hb.done
}
