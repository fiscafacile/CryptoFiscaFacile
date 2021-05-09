package binance

import (
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Binance struct {
	api           api
	csvTXs        []csvTX
	TXsByCategory wallet.TXsByCategory
	done          chan error
	Sources       source.Sources
}

func New() *Binance {
	b := &Binance{}
	b.done = make(chan error)
	b.TXsByCategory = make(wallet.TXsByCategory)
	b.Sources = make(source.Sources)
	return b
}

func (b *Binance) GetAPIExchangeTXs(loc *time.Location) {
	err := b.api.getAllTXs(loc)
	if err != nil {
		b.done <- err
		return
	}
	b.TXsByCategory.Add(b.api.txsByCategory)
	b.done <- nil
}

func (b *Binance) WaitFinish() error {
	return <-b.done
}
