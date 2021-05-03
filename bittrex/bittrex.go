package bittrex

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bittrex struct {
	api           api
	TXsByCategory wallet.TXsByCategory
	transferDone  chan error
	tradesDone    chan error
}

func New() *Bittrex {
	btrx := &Bittrex{}
	btrx.TXsByCategory = make(map[string]wallet.TXs)
	btrx.transferDone = make(chan error)
	btrx.tradesDone = make(chan error)
	return btrx
}

func (btrx *Bittrex) WaitTransfersFinish() error {
	return <-btrx.transferDone
}
func (btrx *Bittrex) WaitTradesFinish() error {
	return <-btrx.tradesDone
}
