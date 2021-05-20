package binance

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"time"
)

type Binance struct {
	api           api
	csvTXs        []csvTX
	done          chan error
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Binance {
	b := &Binance{}
	b.done = make(chan error)
	b.TXsByCategory = make(wallet.TXsByCategory)
	b.Sources = make(source.Sources)
	return b
}

func (b *Binance) GetAPIAllTXs(loc *time.Location) {
	err := b.api.getAllTXs(loc)
	if err != nil {
		b.done <- err
		return
	}
	b.done <- nil
}

func (b *Binance) WaitFinish(account string) error {
	err := <-b.done
	// Merge TX without Duplicates
	b.TXsByCategory.AddUniq(b.api.txsByCategory)
	// Add 3916 Source infos
	if _, ok := b.Sources["Binance"]; ok {
		if b.Sources["Binance"].OpeningDate.After(b.api.firstTimeUsed) {
			src := b.Sources["Binance"]
			src.OpeningDate = b.api.firstTimeUsed
			b.Sources["Binance"] = src
		}
		if b.Sources["Binance"].ClosingDate.Before(b.api.lastTimeUsed) {
			src := b.Sources["Binance"]
			src.ClosingDate = b.api.lastTimeUsed
			b.Sources["Binance"] = src
		}
	} else {
		b.Sources["Binance"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   b.api.firstTimeUsed,
			ClosingDate:   b.api.lastTimeUsed,
			LegalName:     "Binance Europe Services Limited",
			Address:       "LEVEL G (OFFICE 1/1235), QUANTUM HOUSE,75 ABATE RIGORD STREET, TA' XBIEXXBX 1120\nMalta",
			URL:           "https://www.binance.com/fr",
		}
	}
	return err
}
