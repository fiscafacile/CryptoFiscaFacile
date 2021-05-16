package etherscan

import (
	"strings"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type address struct {
	address     string
	description string
}

type Etherscan struct {
	api           api
	addresses     []address
	done          chan error
	TXsByCategory wallet.TXsByCategory
}

func New() *Etherscan {
	ethsc := &Etherscan{}
	ethsc.done = make(chan error)
	ethsc.TXsByCategory = make(map[string]wallet.TXs)
	return ethsc
}

func (ethsc *Etherscan) AddListAddresses(list []string) {
	for _, add := range list {
		ethsc.addresses = append(ethsc.addresses, address{address: strings.ToLower(add)})
	}
}

func (ethsc *Etherscan) GetAPITXs(cat category.Category) {
	addresses := []string{}
	for _, a := range ethsc.addresses {
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
