package localbitcoin

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type LocalBitcoin struct {
	CsvTXsTrade    []CsvTXTrade
	CsvTXsTransfer []CsvTXTransfer
	TXsByCategory  wallet.TXsByCategory
}

func New() *LocalBitcoin {
	lb := &LocalBitcoin{}
	lb.TXsByCategory = make(map[string]wallet.TXs)
	return lb
}
