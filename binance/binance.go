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
	b.TXsByCategory.Add(b.api.txsByCategory)
	b.done <- nil
}

func (b *Binance) WaitFinish() error {
	return <-b.done
}
