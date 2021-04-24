package etherscan

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Etherscan struct {
	api            API
	csvAddresses   []csvAddress
	apiNormalTXs   []apiNormalTX
	apiInternalTXs []apiInternalTX
	apiERC20TXs    []apiERC20TX
	done           chan error
	TXsByCategory  wallet.TXsByCategory
}

func New() *Etherscan {
	ethsc := &Etherscan{}
	ethsc.TXsByCategory = make(map[string]wallet.TXs)
	ethsc.done = make(chan error)
	return ethsc
}

func (ethsc *Etherscan) ownAddress(add string) bool {
	for _, a := range ethsc.csvAddresses {
		if a.address == add {
			return true
		}
	}
	return false
}

func (ethsc *Etherscan) WaitFinish() error {
	return <-ethsc.done
}
