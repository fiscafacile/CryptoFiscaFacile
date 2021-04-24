package bluewallet

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type BlueWallet struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
	Wallets       wallet.Wallets
}

func New() *BlueWallet {
	bw := &BlueWallet{}
	bw.TXsByCategory = make(map[string]wallet.TXs)
	bw.Wallets = make(map[string]decimal.Decimal)
	return bw
}
