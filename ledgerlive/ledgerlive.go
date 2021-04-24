package ledgerlive

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type LedgerLive struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *LedgerLive {
	ll := &LedgerLive{}
	ll.TXsByCategory = make(map[string]wallet.TXs)
	return ll
}
