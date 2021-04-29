package blockstream

import (
	"fmt"
	"log"
	"time"

	"github.com/anaskhan96/base58check"
	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

func (blkst *Blockstream) DetectBTG(b *btc.BTC) {
	btgForkDate := time.Date(2017, time.August, 1, 15, 16, 0, 0, time.UTC)
	w := b.TXsByCategory.GetWallets(btgForkDate, false, false)
	w.Println("BTC (at time of BTG Fork)", "BTC")
	fmt.Println("Addresses :")
	for _, a := range b.CSVAddresses {
		bal, err := blkst.GetAddressBalanceAtDate(a.Address, btgForkDate)
		if err != nil {
			log.Println("")
			break
		}
		if !bal.IsZero() {
			decoded, err := base58check.Decode(a.Address)
			if err != nil {
				log.Println("BTG base58 Decode error", a.Address, err)
			} else {
				version := "26"
				if a.Address[0] == '3' {
					version = "17"
				}
				encoded, err := base58check.Encode(version, decoded[2:])
				if err != nil {
					log.Println("BTG base58 Encode error", decoded, err)
				} else {
					fmt.Println("  ", encoded, "balance", bal)
					t := wallet.TX{Timestamp: btgForkDate, Note: "Blockstream API : 491407 Bitcoin Gold Fork from " + a.Address + " to " + encoded}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BTG", Amount: bal})
					b.TXsByCategory["Forks"] = append(b.TXsByCategory["Forks"], t)
				}
			}
		}
	}
}
