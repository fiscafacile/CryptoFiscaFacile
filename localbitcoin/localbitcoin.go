package localbitcoin

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type LocalBitcoin struct {
	CsvTXsTrade    []CsvTXTrade
	CsvTXsTransfer []CsvTXTransfer
	Accounts       wallet.Accounts
}

func New() *LocalBitcoin {
	lb := &LocalBitcoin{}
	lb.Accounts = make(map[string]wallet.TXs)
	return lb
}
