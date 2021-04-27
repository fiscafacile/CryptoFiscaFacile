package blockchain

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type BlockChain struct {
	jsonTXs       []JsonTX
	TXsByCategory wallet.TXsByCategory
}

func New() *BlockChain {
	cb := &BlockChain{}
	cb.TXsByCategory = make(map[string]wallet.TXs)
	return cb
}
