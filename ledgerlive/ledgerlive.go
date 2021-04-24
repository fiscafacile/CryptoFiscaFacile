package ledgerlive

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type LedgerLive struct {
	CsvTXs   []CsvTX
	Accounts wallet.Accounts
}

func New() *LedgerLive {
	ll := &LedgerLive{}
	ll.Accounts = make(map[string]wallet.TXs)
	return ll
}
