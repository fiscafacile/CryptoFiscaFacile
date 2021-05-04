package etherscan

import (
	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Etherscan struct {
	api           api
	csvAddresses  []csvAddress
	done          chan error
	TXsByCategory wallet.TXsByCategory
}

func New() *Etherscan {
	ethsc := &Etherscan{}
	ethsc.done = make(chan error)
	ethsc.TXsByCategory = make(map[string]wallet.TXs)
	return ethsc
}

func (ethsc *Etherscan) GetAPITXs(cat category.Category) {
	addresses := []string{}
	for _, a := range ethsc.csvAddresses {
		addresses = append(addresses, a.address)
	}
	err := ethsc.api.getAllTXs(addresses, cat)
	if err != nil {
		ethsc.done <- err
		return
	}
	ethsc.TXsByCategory.Add(ethsc.api.txsByCategory)
	ethsc.done <- nil
}

func (ethsc *Etherscan) WaitFinish() error {
	return <-ethsc.done
}
