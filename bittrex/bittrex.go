package bittrex

import (
	"sync"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Bittrex struct {
	api           api
	done          chan error
	mutex         sync.Mutex
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Bittrex {
	btrx := &Bittrex{}
	btrx.done = make(chan error)
	btrx.TXsByCategory = make(wallet.TXsByCategory)
	btrx.Sources = make(source.Sources)
	return btrx
}

func (btrx *Bittrex) GetAPIAllTXs() {
	err := btrx.api.getAllTXs()
	if err != nil {
		btrx.done <- err
		return
	}
	btrx.mutex.Lock()
	for k, v := range btrx.api.txsByCategory {
		if k == "Exchanges" {
			for _, tx := range v {
				found := false
				for _, t := range btrx.TXsByCategory["Exchanges"] {
					if t.ID == tx.ID {
						found = true
						break
					}
				}
				if !found {
					btrx.TXsByCategory[k] = append(btrx.TXsByCategory[k], tx)
				}
			}
		} else {
			btrx.TXsByCategory[k] = append(btrx.TXsByCategory[k], v...)
		}
	}
	btrx.mutex.Unlock()
	if _, ok := btrx.Sources["Bittrex"]; ok {
		if btrx.Sources["Bittrex"].OpeningDate.After(btrx.api.firstTimeUsed) {
			src := btrx.Sources["Bittrex"]
			src.OpeningDate = btrx.api.firstTimeUsed
			btrx.Sources["Bittrex"] = src
		}
		if btrx.Sources["Bittrex"].ClosingDate.Before(btrx.api.lastTimeUsed) {
			src := btrx.Sources["Bittrex"]
			src.ClosingDate = btrx.api.lastTimeUsed
			btrx.Sources["Bittrex"] = src
		}
	} else {
		btrx.Sources["Bittrex"] = source.Source{
			Crypto:        true,
			AccountNumber: "emailAROBASEdomainPOINTcom",
			OpeningDate:   btrx.api.firstTimeUsed,
			ClosingDate:   btrx.api.lastTimeUsed,
			LegalName:     "Bittrex International GmbH",
			Address:       "Dr. Grass-Strasse 12, 9490 Vaduz,\nPrincipality of Liechtenstein",
			URL:           "https://global.bittrex.com",
		}
	}
	btrx.done <- nil
}

func (btrx *Bittrex) WaitFinish() error {
	return <-btrx.done
}
