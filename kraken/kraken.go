package kraken

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Kraken struct {
	csvTXs        []csvTX
	TXsByCategory wallet.TXsByCategory
	transferDone  chan error
	tradesDone    chan error
}

func New() *Kraken {
	kr := &Kraken{}
	kr.TXsByCategory = make(map[string]wallet.TXs)
	kr.transferDone = make(chan error)
	kr.tradesDone = make(chan error)
	return kr
}

// func (kr *Kraken) WaitTransfersFinish() error {
// 	return <-kr.transferDone
// }
func (kr *Kraken) WaitTradesFinish() error {
	return <-kr.tradesDone
}
