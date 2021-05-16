package blockstream

import (
	"fmt"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

func (blkst *Blockstream) DetectBCD(b *btc.BTC) {
	bcdForkDate := time.Date(2017, time.November, 24, 10, 20, 0, 0, time.UTC)
	w := b.TXsByCategory.GetWallets(bcdForkDate, false, false)
	w.Println("BTC (at time of BCD Fork)", "BTC")
	fmt.Println("Addresses :")
	for _, a := range b.Addresses {
		bal, err := blkst.GetAddressBalanceAtDate(a.Address, bcdForkDate)
		if err != nil {
			log.Println("")
			break
		}
		if !bal.IsZero() {
			bal = bal.Mul(decimal.NewFromInt(int64(10)))
			fmt.Println("  ", a.Address, "balance", bal)
			t := wallet.TX{Timestamp: bcdForkDate, Note: "Blockstream API : 495866 Bitcoin Diamond Fork on " + a.Address}
			t.Items = make(map[string]wallet.Currencies)
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BCD", Amount: bal})
			b.TXsByCategory["Forks"] = append(b.TXsByCategory["Forks"], t)
		}
	}
}
