package wallet

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
)

type Currency struct {
	Code   string
	Amount decimal.Decimal
}

type WalletCurrencies map[string]decimal.Decimal

type Wallets struct {
	Date       time.Time
	Currencies WalletCurrencies
}

type TX struct {
	Timestamp time.Time
	Items     map[string][]Currency
	Note      string
}

type TXs []TX

type TXsByCategory map[string]TXs

func (c *Currency) IsFiat() bool {
	if c.Code == "EUR" ||
		c.Code == "USD" ||
		c.Code == "HKD" {
		return true
	}
	return false
}

func (c *Currency) Println() {
	fmt.Println(c.Amount, c.Code)
}

func (c Currency) GetExchangeRate(date time.Time, to string) (rate decimal.Decimal, err error) {
	gecko, err := NewCoinGeckoAPI()
	if err == nil {
		ratesCG, err := gecko.GetExchangeRates(date, c.Code)
		if err == nil {
			// log.Println("ratesCG : ", ratesCG)
			for _, r := range ratesCG.Rates {
				if r.Quote == to {
					return r.Rate, nil
				}
			}
			// } else {
			// 	log.Println("gecko.GetExchangeRates :", err)
		}
		// } else {
		// 	log.Println("NewCoinGeckoAPI :", err)
	}
	var api CoinAPI
	rates, err := api.GetExchangeRates(date, to)
	if err == nil {
		for _, r := range rates.Rates {
			if r.Quote == c.Code {
				return r.Rate, nil
			}
		}
		// } else {
		// 	log.Println("CoinAPI.GetExchangeRates :", err)
	}
	return rate, errors.New("Cannot find rate for " + to + c.Code + " at " + date.String())
}

func (w Wallets) Round() {
	for k, v := range w.Currencies {
		if k == "BAB" {
			if v.Abs().LessThan(decimal.NewFromInt(1)) {
				delete(w.Currencies, k)
			}
		} else if k == "EUR" || k == "USD" || k == "HKD" || k == "CRO" || k == "USDC" || k == "USDT" || k == "sUSD" || k == "XRP" || k == "IOT" {
			if v.Abs().LessThan(decimal.NewFromFloat(0.5)) {
				delete(w.Currencies, k)
			}
		} else if k == "LPT" {
			if v.Abs().LessThan(decimal.NewFromInt(100)) {
				delete(w.Currencies, k)
			}
		} else {
			if v.Abs().LessThan(decimal.NewFromFloat(0.01)) {
				delete(w.Currencies, k)
			}
		}
	}
}

func (w Wallets) CalculateTotalValue(native string) (totalValue Currency, err error) {
	totalValue.Code = native
	for k, v := range w.Currencies {
		if k == native {
			totalValue.Amount = totalValue.Amount.Add(v)
		} else {
			c := Currency{Code: k, Amount: v}
			rate, err := c.GetExchangeRate(w.Date, native)
			if err != nil {
				log.Println("Cannot find rate for", k, "at", w.Date)
			} else {
				totalValue.Amount = totalValue.Amount.Add(rate.Mul(v))
			}
		}
	}
	return
}

func (w Wallets) Println(name string) {
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("| " + name + " Wallet |")
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("Time :", w.Date.UTC())
	fmt.Println("Amounts :")
	keys := make([]string, 0, len(w.Currencies))
	for k := range w.Currencies {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println("  ", w.Currencies[k], k)
	}
}

func (tx *TX) SimilarDate(delta time.Duration, t time.Time) bool {
	if tx.Timestamp.Sub(t) < delta &&
		t.Sub(tx.Timestamp) < delta {
		return true
	}
	return false
}

func (tx TX) Println() {
	fmt.Println("Time :", tx.Timestamp.UTC())
	for k, v := range tx.Items {
		fmt.Println(k, ":")
		for _, i := range v {
			fmt.Println("  ", i.Amount.String(), i.Code)
		}
	}
	fmt.Println("Note :", tx.Note)
}

func (txs TXs) SortByDate(chrono bool) {
	if chrono {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Timestamp.Before(txs[j].Timestamp)
		})
	} else {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Timestamp.After(txs[j].Timestamp)
		})
	}
}

func (acc TXsByCategory) GetWallets(date time.Time, includeFiat bool) (w Wallets) {
	w.Date = date
	w.Currencies = make(WalletCurrencies)
	for _, a := range acc {
		for _, tx := range a {
			if tx.Timestamp.Before(date) {
				for k, i := range tx.Items {
					for _, c := range i {
						if c.Code != "" &&
							(includeFiat || !c.IsFiat()) {
							if k == "Fee" || k == "From" {
								w.Currencies[c.Code] = w.Currencies[c.Code].Sub(c.Amount)
							} else {
								w.Currencies[c.Code] = w.Currencies[c.Code].Add(c.Amount)
							}
						}
					}
				}
			}
		}
	}
	w.Round()
	return
}

func (acc TXsByCategory) Add(a TXsByCategory) {
	for k, v := range a {
		acc[k] = append(acc[k], v...)
	}
}

func (acc TXsByCategory) FindTransfers() TXsByCategory {
	var realDeposits TXs
	var realWithdrawals TXs
	similarTimeDelta := 12 * time.Hour
	for _, depTX := range acc["Deposits"] {
		found := false
		depFees := decimal.NewFromInt(0)
		if _, ok := depTX.Items["Fee"]; ok {
			if len(depTX.Items["Fee"]) == 1 {
				depFees = depTX.Items["Fee"][0].Amount
			}
		}
		for _, witTX := range acc["Withdrawals"] {
			if depTX.Items["To"][0].Code == witTX.Items["From"][0].Code &&
				depTX.SimilarDate(similarTimeDelta, witTX.Timestamp) &&
				strings.Split(depTX.Note, ":")[0] != strings.Split(witTX.Note, ":")[0] {
				witFees := decimal.NewFromInt(0)
				if _, ok := witTX.Items["Fee"]; ok {
					if len(witTX.Items["Fee"]) == 1 {
						witFees = witTX.Items["Fee"][0].Amount
					}
				}
				if depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount) ||
					depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(witFees)) ||
					depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(depFees)) {
					found = true
					t := TX{Timestamp: witTX.Timestamp, Note: witTX.Note + " => " + depTX.Note}
					t.Items = make(map[string][]Currency)
					t.Items["To"] = append(t.Items["To"], depTX.Items["To"]...)
					t.Items["From"] = append(t.Items["From"], witTX.Items["From"]...)
					if _, ok := witTX.Items["Fee"]; ok {
						t.Items["Fee"] = append(t.Items["Fee"], witTX.Items["Fee"]...)
					}
					if _, ok := depTX.Items["Fee"]; ok {
						for _, df := range depTX.Items["Fee"] {
							missing := true
							for _, f := range t.Items["Fee"] {
								if f.Code == df.Code &&
									f.Amount.Equal(df.Amount) {
									missing = false
								}
							}
							if missing {
								t.Items["Fee"] = append(t.Items["Fee"], df)
							}
						}
					}
					acc["Transfers"] = append(acc["Transfers"], t)
					break
					// } else {
					// 	spew.Dump(depTX)
					// 	spew.Dump(witTX)
				}
			}
		}
		if !found {
			realDeposits = append(realDeposits, depTX)
		}
	}
	for _, witTX := range acc["Withdrawals"] {
		found := false
		witFees := decimal.NewFromInt(0)
		if _, ok := witTX.Items["Fee"]; ok {
			if len(witTX.Items["Fee"]) == 1 {
				witFees = witTX.Items["Fee"][0].Amount
			}
		}
		for _, depTX := range acc["Deposits"] {
			depFees := decimal.NewFromInt(0)
			if _, ok := depTX.Items["Fee"]; ok {
				if len(depTX.Items["Fee"]) == 1 {
					depFees = depTX.Items["Fee"][0].Amount
				}
			}
			if depTX.Items["To"][0].Code == witTX.Items["From"][0].Code &&
				depTX.SimilarDate(similarTimeDelta, witTX.Timestamp) &&
				strings.Split(depTX.Note, ":")[0] != strings.Split(witTX.Note, ":")[0] {
				if depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount) ||
					depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(witFees)) ||
					depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(depFees)) {
					found = true
					break
				}
			}
		}
		if !found {
			realWithdrawals = append(realWithdrawals, witTX)
		}
	}
	acc["Deposits"] = realDeposits
	acc["Withdrawals"] = realWithdrawals
	return acc
}

func (acc TXsByCategory) FindCashInOut() TXsByCategory {
	var realExchanges TXs
	for _, exTX := range acc["Exchanges"] {
		fromHasFiat := false
		for _, i := range exTX.Items["From"] {
			if i.IsFiat() {
				fromHasFiat = true
			}
		}
		toHasFiat := false
		for _, i := range exTX.Items["To"] {
			if i.IsFiat() {
				toHasFiat = true
			}
		}
		if fromHasFiat {
			acc["CashIn"] = append(acc["CashIn"], exTX)
		} else if toHasFiat {
			acc["CashOut"] = append(acc["CashOut"], exTX)
		} else {
			realExchanges = append(realExchanges, exTX)
		}
	}
	if len(realExchanges) > 0 {
		acc["Exchanges"] = realExchanges
	}
	// var realDeposits TXs
	// for _, depTX := range acc["Deposits"] {
	// 	toHasFiat := false
	// 	for _, i := range depTX.Items["To"] {
	// 		if i.IsFiat() {
	// 			toHasFiat = true
	// 		}
	// 	}
	// 	if toHasFiat {
	// 		acc["CashIn"] = append(acc["CashIn"], depTX)
	// 	} else {
	// 		realDeposits = append(realDeposits, depTX)
	// 	}
	// }
	// acc["Deposits"] = realDeposits
	var realWithdrawals TXs
	for _, witTX := range acc["Withdrawals"] {
		fromHasFiat := false
		// for _, i := range witTX.Items["From"] {
		// 	if i.IsFiat() {
		// 		fromHasFiat = true
		// 	}
		// }
		if fromHasFiat {
			acc["CashOut"] = append(acc["CashOut"], witTX)
		} else {
			realWithdrawals = append(realWithdrawals, witTX)
		}
	}
	if len(realWithdrawals) > 0 {
		acc["Withdrawals"] = realWithdrawals
	}
	return acc
}

func (acc TXsByCategory) SortTXsByDate(chrono bool) {
	for k := range acc {
		acc[k].SortByDate(chrono)
	}
}

func (acc TXsByCategory) PrintStats() {
	fmt.Println("---------------------------")
	fmt.Println("| List of TXs By Category |")
	fmt.Println("---------------------------")
	keys := make([]string, 0, len(acc))
	for k := range acc {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, ":", len(acc[k]), "TXs")
	}
}

func (acc TXsByCategory) PrintUnjustifiedWithdrawals(loc *time.Location) {
	have := false
	for _, tx := range acc["Withdrawals"] {
		if tx.Timestamp.After(time.Date(2018, time.December, 31, 23, 59, 59, 999, loc)) {
			if !have {
				fmt.Println("-----------------------------------")
				fmt.Println("| List of Unjustified Withdrawals |")
			}
			fmt.Println("-----------------------------------")
			tx.Println()
			have = true
		}
	}
	if !have {
		fmt.Println("All Withdrawals are justified")
	}
}
