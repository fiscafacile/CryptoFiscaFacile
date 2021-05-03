package kraken

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Kraken struct {
	// api           api
	TXsByCategory wallet.TXsByCategory
	transferDone  chan error
	tradesDone    chan error
}

func New() *Kraken {
	krkn := &Kraken{}
	krkn.TXsByCategory = make(map[string]wallet.TXs)
	krkn.transferDone = make(chan error)
	krkn.tradesDone = make(chan error)
	return krkn
}

// func (krkn *Kraken) WaitTransfersFinish() error {
// 	return <-krkn.transferDone
// }
func (krkn *Kraken) WaitTradesFinish() error {
	return <-krkn.tradesDone
}
