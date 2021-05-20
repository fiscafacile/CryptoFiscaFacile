package hitbtc

import (
	"strings"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
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
	hb.done <- nil
}

func (hb *HitBTC) WaitFinish(account string) error {
	err := <-hb.done
	// Merge TX without Duplicates
	hb.TXsByCategory.AddUniq(hb.api.txsByCategory)
	// Add 3916 Source infos
	if _, ok := hb.Sources["HitBTC"]; ok {
		if hb.Sources["HitBTC"].OpeningDate.After(hb.api.firstTimeUsed) {
			src := hb.Sources["HitBTC"]
			src.OpeningDate = hb.api.firstTimeUsed
			hb.Sources["HitBTC"] = src
		}
		if hb.Sources["HitBTC"].ClosingDate.Before(hb.api.lastTimeUsed) {
			src := hb.Sources["HitBTC"]
			src.ClosingDate = hb.api.lastTimeUsed
			hb.Sources["HitBTC"] = src
		}
	} else {
		hb.Sources["HitBTC"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   hb.api.firstTimeUsed,
			ClosingDate:   hb.api.lastTimeUsed,
			LegalName:     "Hit Tech Solutions Development Ltd.",
			Address:       "Suite 15, Oliaji Trade Centre, Francis Rachel Street,\nVictoria, Mahe,\nSeychelles",
			URL:           "https://hitbtc.com",
		}
	}
	return err
}

func csvCurrencyCure(c string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(c, "BCHSV", "BSV"), "BCHABC", "BCH"), "BCCF", "BCH"))
}

func apiCurrencyCure(c string) string {
	// https://blog.hitbtc.com/we-will-change-the-ticker-of-bchabc-to-bch-and-bchsv-will-be-displayed-as-hbv/
	return strings.ReplaceAll(strings.ReplaceAll(c, "BCHA", "BCH"), "BCHOLD", "BCH")
}
