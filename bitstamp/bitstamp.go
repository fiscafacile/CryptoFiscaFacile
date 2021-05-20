package bitstamp

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bitstamp struct {
	api           api
	csvTXs        []csvTX
	done          chan error
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Bitstamp {
	bs := &Bitstamp{}
	bs.done = make(chan error)
	bs.TXsByCategory = make(wallet.TXsByCategory)
	bs.Sources = make(map[string]source.Source)
	return bs
}

func (bs *Bitstamp) GetAPIAllTXs() {
	err := bs.api.getAllTXs()
	if err != nil {
		bs.done <- err
		return
	}
	bs.done <- nil
}

func (bs *Bitstamp) WaitFinish(account string) error {
	err := <-bs.done
	// Merge TX without Duplicates
	bs.TXsByCategory.AddUniq(bs.api.txsByCategory)
	// Add 3916 Source infos
	if _, ok := bs.Sources["Bitstamp"]; ok {
		if bs.Sources["Bitstamp"].OpeningDate.After(bs.api.firstTimeUsed) {
			src := bs.Sources["Bitstamp"]
			src.OpeningDate = bs.api.firstTimeUsed
			bs.Sources["Bitstamp"] = src
		}
		if bs.Sources["Bitstamp"].ClosingDate.Before(bs.api.lastTimeUsed) {
			src := bs.Sources["Bitstamp"]
			src.ClosingDate = bs.api.lastTimeUsed
			bs.Sources["Bitstamp"] = src
		}
	} else {
		bs.Sources["Bitstamp"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   bs.api.firstTimeUsed,
			ClosingDate:   bs.api.lastTimeUsed,
			LegalName:     "Bitstamp Ltd",
			Address:       "5 New Street Square,\nLondon EC4A 3TW,\nRoyaume-Uni",
			URL:           "https://bitstamp.com",
		}
	}
	return err
}
