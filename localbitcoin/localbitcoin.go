package localbitcoin

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type LocalBitcoin struct {
	CsvTXsTrade    []CsvTXTrade
	CsvTXsTransfer []CsvTXTransfer
	TXsByCategory  wallet.TXsByCategory
	Sources        source.Sources
}

func New() *LocalBitcoin {
	lb := &LocalBitcoin{}
	lb.TXsByCategory = make(map[string]wallet.TXs)
	lb.Sources = make(source.Sources)
	return lb
}
