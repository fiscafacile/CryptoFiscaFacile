package bitstamp

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/utils"
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
	bs.TXsByCategory.Add(bs.api.txsByCategory)
	if _, ok := bs.Sources["Bitstamp"]; !ok {
		bs.Sources["Bitstamp"] = source.Source{
			Crypto:        true,
			AccountNumber: utils.RemoveSymbol("email@domain.com"),
			OpeningDate:   bs.api.firstTimeUsed,
			ClosingDate:   bs.api.lastTimeUsed,
			LegalName:     "Bitstamp Ltd",
			Address:       "5 New Street Square,\nLondon EC4A 3TW,\nRoyaume-Uni",
			URL:           "https://bitstamp.com",
		}
	}
	bs.done <- nil
}

func (bs *Bitstamp) WaitFinish() error {
	return <-bs.done
}
