package kraken

import (
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Kraken struct {
	api           api
	csvTXs        []csvTX
	done          chan error
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Kraken {
	kr := &Kraken{}
	kr.done = make(chan error)
	kr.TXsByCategory = make(map[string]wallet.TXs)
	kr.Sources = make(source.Sources)
	return kr
}

func (kr *Kraken) GetAPIAllTXs() {
	err := kr.api.getAPITxs()
	if err != nil {
		kr.done <- err
		return
	}
	kr.TXsByCategory.Add(kr.api.txsByCategory)
	kr.Sources["Kraken"] = source.Source{
		Crypto:        true,
		AccountNumber: "emailAROBASEdomainPOINTcom",
		OpeningDate:   kr.api.firstTimeUsed,
		ClosingDate:   kr.api.lastTimeUsed,
		LegalName:     "Payward Ltd.",
		Address:       "6th Floor,\nOne London Wall,\nLondon, EC2Y 5EB,\nRoyaume-Uni",
		URL:           "https://www.kraken.com",
	}
	kr.done <- nil
}

func (kr *Kraken) WaitFinish() error {
	return <-kr.done
}
