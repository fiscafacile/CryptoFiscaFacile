package wallet

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
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

type Currencies []Currency

type WalletCurrencies map[string]decimal.Decimal

type Wallets struct {
	Date       time.Time
	Currencies WalletCurrencies
}

type TX struct {
	Timestamp time.Time
	Items     map[string]Currencies
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

func (c *Currency) Println(filter string) {
	if strings.Contains(filter, c.Code) ||
		filter == "" {
		fmt.Println(c.Amount, c.Code)
	}
}

func (c Currency) GetExchangeRate(date time.Time, to string) (rate decimal.Decimal, err error) {
	if !c.IsFiat() {
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
	}
	var layer CoinLayer
	ratesCL, err := layer.GetExchangeRates(date, to)
	if err == nil {
		return decimal.NewFromFloat(ratesCL.Rates[c.Code]), nil
		// } else {
		// 	log.Println("CoinLayer.GetExchangeRates :", err)
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

func (wc WalletCurrencies) Add(a WalletCurrencies) {
	for k, v := range a {
		wc[k] = wc[k].Add(v)
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

func (w Wallets) Round(rounding bool) {
	for k, v := range w.Currencies {
		if rounding {
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
		} else {
			if v.IsZero() {
				delete(w.Currencies, k)
			}
		}
	}
}

func (w Wallets) Println(name string, filter string) {
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("| " + name + " Wallet |")
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("Time :", w.Date.UTC())
	fmt.Println("Amounts :")
	keys := make([]string, 0, len(w.Currencies))
	for k := range w.Currencies {
		if filter == "" || strings.Contains(filter, k) {
			keys = append(keys, k)
		}
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

func (tx TX) Base64String(anonymous bool) string {
	var temp bytes.Buffer
	binary.Write(&temp, binary.LittleEndian, tx)
	return base64.StdEncoding.EncodeToString(append([]byte("TX"), temp.Bytes()...))
}

func (tx TX) Println(filter string) (printed bool) {
	if filter == "" {
		printed = true
	} else {
		printed = false
	}
	toPrint := fmt.Sprintln("Time :", tx.Timestamp.UTC())
	for k, v := range tx.Items {
		toPrint += fmt.Sprintln(k, ":")
		for _, i := range v {
			toPrint += fmt.Sprintln("  ", i.Amount.String(), i.Code)
			if strings.Contains(filter, i.Code) {
				printed = true
			}
		}
	}
	toPrint += fmt.Sprintln("Note :", tx.Note)
	if printed {
		fmt.Print(toPrint)
	}
	return
}

func (tx TX) GetCurrencies(includeFiat, includeFee bool) (cs WalletCurrencies) {
	cs = make(WalletCurrencies)
	for k, i := range tx.Items {
		for _, c := range i {
			if c.Code != "" &&
				(includeFiat || !c.IsFiat()) {
				if (k == "Fee" && includeFee) || k == "From" {
					cs[c.Code] = cs[c.Code].Sub(c.Amount)
				} else if k == "To" {
					cs[c.Code] = cs[c.Code].Add(c.Amount)
				}
			}
		}
	}
	return
}

func (txs TXs) Println(name string, filter string) {
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("| TXs in " + name + " |")
	fmt.Println(strings.Repeat("-", len(name)+11))
	printed := false
	for _, tx := range txs {
		if printed {
			fmt.Println(strings.Repeat("-", len(name)+11))
		}
		printed = tx.Println(filter)
	}
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

func (txs TXsByCategory) Println(filter string) {
	for k, v := range txs {
		v.Println("Category "+k, filter)
	}
}

func (txs TXsByCategory) GetWallets(date time.Time, includeFiat bool, rounding bool) (w Wallets) {
	w.Date = date
	w.Currencies = make(WalletCurrencies)
	for _, a := range txs {
		for _, tx := range a {
			if tx.Timestamp.Before(date) {
				txcs := tx.GetCurrencies(includeFiat, true)
				w.Currencies.Add(txcs)
			}
		}
	}
	w.Round(rounding)
	return
}

func (txs TXsByCategory) Add(a TXsByCategory) {
	for k, v := range a {
		txs[k] = append(txs[k], v...)
	}
}

func (txs TXsByCategory) FindTransfers() TXsByCategory {
	var realDeposits TXs
	var realWithdrawals TXs
	similarTimeDelta := 12 * time.Hour
	for _, depTX := range txs["Deposits"] {
		found := false
		depFees := decimal.NewFromInt(0)
		if _, ok := depTX.Items["Fee"]; ok {
			for _, f := range depTX.Items["Fee"] {
				depFees = depFees.Add(f.Amount)
			}
		}
		for _, witTX := range txs["Withdrawals"] {
			if depTX.Items["To"][0].Code == witTX.Items["From"][0].Code &&
				depTX.SimilarDate(similarTimeDelta, witTX.Timestamp) &&
				strings.Split(depTX.Note, ":")[0] != strings.Split(witTX.Note, ":")[0] {
				witFees := decimal.NewFromInt(0)
				if _, ok := witTX.Items["Fee"]; ok {
					for _, f := range witTX.Items["Fee"] {
						witFees = witFees.Add(f.Amount)
					}
				}
				// log.Println("Here")
				// depTX.Println("")
				// witTX.Println("")
				if depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount) ||
					depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(witFees)) ||
					depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(depFees)) {
					found = true
					t := TX{Timestamp: witTX.Timestamp, Note: witTX.Note + " => " + depTX.Note}
					t.Items = make(map[string]Currencies)
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
					txs["Transfers"] = append(txs["Transfers"], t)
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
	for _, witTX := range txs["Withdrawals"] {
		found := false
		witFees := decimal.NewFromInt(0)
		if _, ok := witTX.Items["Fee"]; ok {
			for _, f := range witTX.Items["Fee"] {
				witFees = witFees.Add(f.Amount)
			}
		}
		for _, depTX := range txs["Deposits"] {
			depFees := decimal.NewFromInt(0)
			if _, ok := depTX.Items["Fee"]; ok {
				for _, f := range depTX.Items["Fee"] {
					depFees = depFees.Add(f.Amount)
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
	txs["Deposits"] = realDeposits
	txs["Withdrawals"] = realWithdrawals
	return txs
}

func (txs TXsByCategory) FindCashInOut(native string) TXsByCategory {
	var realExchanges TXs
	for _, exTX := range txs["Exchanges"] {
		fromHasFiat := false
		for _, i := range exTX.Items["From"] {
			if i.IsFiat() {
				fromHasFiat = true
			}
		}
		toHasFiat := false
		for _, c := range exTX.Items["To"] {
			if c.IsFiat() {
				toHasFiat = true
			}
		}
		if fromHasFiat {
			txs["CashIn"] = append(txs["CashIn"], exTX)
		} else if toHasFiat {
			txs["CashOut"] = append(txs["CashOut"], exTX)
		} else {
			realExchanges = append(realExchanges, exTX)
		}
	}
	if len(realExchanges) > 0 {
		txs["Exchanges"] = realExchanges
	} else {
		delete(txs, "Exchanges")
	}
	// Integrate CashBacks into CashIn with Reversal cancelation
	var realCashBacks TXs
	var toCancel Currencies
	txs["Cashbacks"].SortByDate(false)
	// txs["Cashbacks"].Println("Cashbacks")
	for _, cbTX := range txs["Cashbacks"] {
		// cbTX.Println()
		isReversal := false
		for _, c := range cbTX.Items["From"] {
			toCancel = append(toCancel, c)
			isReversal = true
			// log.Println("isReversal, toCancel", toCancel)
		}
		if !isReversal {
			// log.Println("!isReversal, toCancel", toCancel)
			for _, c := range cbTX.Items["To"] {
				canceled := false
				// c.Println()
				var toCancelNew Currencies
				for _, tc := range toCancel {
					if c.Amount.Equal(tc.Amount) {
						// log.Println("!isReversal, canceled", tc)
						canceled = true
					} else {
						toCancelNew = append(toCancelNew, tc)
					}
				}
				// log.Println("toCancelNew", toCancelNew)
				if canceled {
					toCancel = toCancelNew
					// log.Println("!isReversal, canceled, toCancel", toCancel)
				} else {
					rate, err := c.GetExchangeRate(cbTX.Timestamp, native)
					if err != nil {
						log.Println(err)
						realCashBacks = append(realCashBacks, cbTX)
					} else {
						cbTX.Items["From"] = append(cbTX.Items["From"], Currency{Code: native, Amount: c.Amount.Mul(rate)})
						txs["CashIn"] = append(txs["CashIn"], cbTX)
					}
				}
			}
		}
	}
	if len(realCashBacks) > 0 {
		txs["Cashbacks"] = realCashBacks
	} else {
		delete(txs, "Cashbacks")
	}
	/*
		var realDeposits TXs
		for _, depTX := range txs["Deposits"] {
			toHasFiat := false
			for _, i := range depTX.Items["To"] {
				if i.IsFiat() {
					toHasFiat = true
				}
			}
			if toHasFiat {
				txs["CashIn"] = append(txs["CashIn"], depTX)
			} else {
				realDeposits = append(realDeposits, depTX)
			}
		}
		txs["Deposits"] = realDeposits
	*/
	/*
		var realWithdrawals TXs
		for _, witTX := range txs["Withdrawals"] {
			fromHasFiat := false
			for _, i := range witTX.Items["From"] {
				if i.IsFiat() {
					fromHasFiat = true
				}
			}
			if fromHasFiat {
				txs["CashOut"] = append(txs["CashOut"], witTX)
				log.Println("Found")
			} else {
				realWithdrawals = append(realWithdrawals, witTX)
			}
		}
		if len(realWithdrawals) > 0 {
			txs["Withdrawals"] = realWithdrawals
		}
	*/
	return txs
}

func (txs TXsByCategory) SortTXsByDate(chrono bool) {
	for k := range txs {
		txs[k].SortByDate(chrono)
	}
}

func (txs TXsByCategory) PrintStats() {
	fmt.Println("-------------------------------")
	fmt.Println("| Quantity of TXs By Category |")
	fmt.Println("-------------------------------")
	keys := make([]string, 0, len(txs))
	for k := range txs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, ":", len(txs[k]), "TXs")
	}
}

func (txs TXsByCategory) CheckConsistency(loc *time.Location) {
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Unjustified Withdrawals (after 2019 Jan 1st) |")
	for _, tx := range txs["Withdrawals"] {
		if tx.Timestamp.After(time.Date(2018, time.December, 31, 23, 59, 59, 999, loc)) {
			fmt.Println("--------------------------------------------------------")
			tx.Println("")
		}
	}
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Non-Zero balance Transfers                   |")
	for _, tx := range txs["Transfers"] {
		txcs := tx.GetCurrencies(false, false)
		for _, v := range txcs {
			if !v.IsZero() {
				fmt.Println("--------------------------------------------------------")
				tx.Println("")
			}
		}
	}
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Negative Amounts TXs                         |")
	for _, v := range txs {
		for _, tx := range v {
			for _, i := range tx.Items {
				for _, c := range i {
					if c.Amount.IsNegative() {
						fmt.Println("--------------------------------------------------------")
						tx.Println("")
					}
				}
			}
		}
	}
	fmt.Println("--------------------------------------------------------")
}
