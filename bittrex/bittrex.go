package bittrex

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bittrex struct {
	api           api
	apiTradeTX    []apiTradeTX
	apiTransferTX []apiTransferTX
	TXsByCategory wallet.TXsByCategory
	done          chan error
}

func New() *Bittrex {
	btrx := &Bittrex{}
	btrx.TXsByCategory = make(map[string]wallet.TXs)
	btrx.done = make(chan error)
	return btrx
}

func (btrx *Bittrex) WaitFinish() error {
	return <-btrx.done
}
