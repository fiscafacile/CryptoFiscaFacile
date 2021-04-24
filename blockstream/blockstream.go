package blockstream

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type Blockstream struct {
	csvPayments  []csvPayment
	csvAddresses []csvAddress
	apiTXs       []apiTX
	done         chan error
	Accounts     wallet.Accounts
}

func New() *Blockstream {
	blkst := &Blockstream{}
	blkst.Accounts = make(map[string]wallet.TXs)
	blkst.done = make(chan error)
	return blkst
}

func (blkst *Blockstream) ownAddress(add string) bool {
	for _, a := range blkst.csvAddresses {
		if a.address == add {
			return true
		}
	}
	return false
}

func (blkst *Blockstream) isTxPayment(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range blkst.csvPayments {
		if a.txID == txid {
			is = true
			desc = a.description
			val = a.value
			curr = a.currency
			return
		}
	}
	return
}

func (blkst *Blockstream) WaitFinish() error {
	return <-blkst.done
}
