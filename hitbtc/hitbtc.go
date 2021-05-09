package hitbtc

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/utils"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type HitBTC struct {
	api               api
	csvTradeTXs       []csvTradeTX
	csvTransactionTXs []csvTransactionTX
	done              chan error
	TXsByCategory     wallet.TXsByCategory
	emails            []string
	Sources           source.Sources
}

func New() *HitBTC {
	hb := &HitBTC{}
	hb.done = make(chan error)
	hb.TXsByCategory = make(wallet.TXsByCategory)
	hb.Sources = make(map[string]source.Source)
	return hb
}

func (hb *HitBTC) GetAPIAllTXs() {
	err := hb.api.getAllTXs()
	if err != nil {
		hb.done <- err
		return
	}
	hb.TXsByCategory.Add(hb.api.txsByCategory)
	if _, ok := hb.Sources["HitBTC"]; !ok {
		hb.Sources["HitBTC"] = source.Source{
			Crypto:        true,
			AccountNumber: utils.RemoveSymbol("email@domain.com"),
			OpeningDate:   hb.api.firstTimeUsed,
			ClosingDate:   hb.api.lastTimeUsed,
			LegalName:     "Hit Tech Solutions Development Ltd.",
			Address:       "Suite 15, Oliaji Trade Centre, Francis Rachel Street,\nVictoria, Mahe,\nSeychelles",
			URL:           "https://hitbtc.com",
		}
	}
	hb.done <- nil
}

func (hb *HitBTC) WaitFinish() error {
	return <-hb.done
}
