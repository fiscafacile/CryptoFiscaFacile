package cryptocom

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvAppCryptoTX struct {
	ID              string
	Timestamp       time.Time
	Description     string
	Currency        string
	Amount          decimal.Decimal
	ToCurrency      string
	ToAmount        decimal.Decimal
	NativeCurrency  string
	NativeAmount    decimal.Decimal
	NativeAmountUSD decimal.Decimal
	Kind            string
}

func (cdc *CryptoCom) ParseCSVAppCrypto(reader io.Reader, cat category.Category, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	hasCashback := false
	firstTimeCashback := time.Now()
	lastTimeCashback := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Crypto.com App CSV Crypto :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Timestamp (UTC)" {
				tx := csvAppCryptoTX{}
				tx.Timestamp, err = time.Parse("2006-01-02 15:04:05", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Timestamp", r[0])
				}
				tx.Description = r[1]
				tx.Currency = r[2]
				tx.Amount, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[3])
				}
				tx.ToCurrency = r[4]
				tx.ToAmount, _ = decimal.NewFromString(r[5])
				tx.NativeCurrency = r[6]
				tx.NativeAmount, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println(SOURCE, "Error Parsing NativeAmount", r[7])
				}
				tx.NativeAmountUSD, err = decimal.NewFromString(r[8])
				if err != nil {
					log.Println(SOURCE, "Error Parsing NativeAmountUSD", r[8])
				}
				tx.Kind = r[9]
				hash := sha256.Sum256([]byte(SOURCE + tx.Timestamp.String()))
				tx.ID = hex.EncodeToString(hash[:])
				cdc.csvAppCryptoTXs = append(cdc.csvAppCryptoTXs, tx)
				if tx.Timestamp.Before(firstTimeUsed) {
					firstTimeUsed = tx.Timestamp
				}
				if tx.Timestamp.After(lastTimeUsed) {
					lastTimeUsed = tx.Timestamp
				}
				// Fill TXsByCategory
				if tx.Kind == "referral_card_cashback" ||
					tx.Kind == "reimbursement" {
					hasCashback = true
					if tx.Timestamp.Before(firstTimeCashback) {
						firstTimeCashback = tx.Timestamp
					}
					if tx.Timestamp.After(lastTimeCashback) {
						lastTimeCashback = tx.Timestamp
					}
				} else {
					if tx.Timestamp.Before(firstTimeUsed) {
						firstTimeUsed = tx.Timestamp
					}
					if tx.Timestamp.After(lastTimeUsed) {
						lastTimeUsed = tx.Timestamp
					}
				}
				if tx.Kind == "dust_conversion_credited" ||
					tx.Kind == "dust_conversion_debited" ||
					tx.Kind == "interest_swap_credited" ||
					tx.Kind == "interest_swap_debited" ||
					tx.Kind == "lockup_swap_credited" ||
					tx.Kind == "lockup_swap_debited" ||
					tx.Kind == "crypto_wallet_swap_credited" ||
					tx.Kind == "crypto_wallet_swap_debited" {
					found := false
					for i, ex := range cdc.TXsByCategory["Exchanges"] {
						if ex.SimilarDate(2*time.Second, tx.Timestamp) &&
							ex.Note[:5] == tx.Kind[:5] {
							found = true
							if tx.Amount.IsPositive() {
								cdc.TXsByCategory["Exchanges"][i].Items["To"] = append(cdc.TXsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
							} else {
								cdc.TXsByCategory["Exchanges"][i].Items["From"] = append(cdc.TXsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							}
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " " + tx.Kind + " " + tx.Description}
						t.Items = make(map[string]wallet.Currencies)
						if tx.Amount.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
							cdc.TXsByCategory["Exchanges"] = append(cdc.TXsByCategory["Exchanges"], t)
						} else {
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							cdc.TXsByCategory["Exchanges"] = append(cdc.TXsByCategory["Exchanges"], t)
						}
					}
				} else if tx.Kind == "crypto_exchange" ||
					tx.Kind == "viban_purchase" {
					t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " " + tx.Kind + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.ToCurrency, Amount: tx.ToAmount})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
					cdc.TXsByCategory["Exchanges"] = append(cdc.TXsByCategory["Exchanges"], t)
				} else if tx.Kind == "card_top_up" {
					t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " " + tx.Kind + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.NativeCurrency, Amount: tx.NativeAmount.Neg()})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
					cdc.TXsByCategory["Exchanges"] = append(cdc.TXsByCategory["Exchanges"], t)
				} else if tx.Kind == "crypto_deposit" ||
					tx.Kind == "viban_deposit" ||
					tx.Kind == "exchange_to_crypto_transfer" ||
					tx.Kind == "admin_wallet_credited" ||
					tx.Kind == "referral_card_cashback" ||
					tx.Kind == "transfer_cashback" ||
					tx.Kind == "reimbursement" ||
					tx.Kind == "crypto_earn_interest_paid" ||
					tx.Kind == "crypto_earn_extra_interest_paid" ||
					tx.Kind == "gift_card_reward" ||
					tx.Kind == "pay_checkout_reward" ||
					tx.Kind == "referral_gift" ||
					tx.Kind == "referral_bonus" ||
					tx.Kind == "mco_stake_reward" ||
					tx.Kind == "supercharger_withdrawal" ||
					tx.Kind == "crypto_purchase" ||
					tx.Kind == "staking_reward" {
					t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " " + tx.Kind + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					if tx.Kind == "crypto_purchase" {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.NativeCurrency, Amount: tx.NativeAmount})
						cdc.TXsByCategory["CashIn"] = append(cdc.TXsByCategory["CashIn"], t)
					} else if tx.Kind == "referral_card_cashback" ||
						tx.Kind == "transfer_cashback" ||
						tx.Kind == "reimbursement" ||
						tx.Kind == "gift_card_reward" ||
						tx.Kind == "referral_gift" ||
						tx.Kind == "pay_checkout_reward" {
						cdc.TXsByCategory["CommercialRebates"] = append(cdc.TXsByCategory["CommercialRebates"], t)
					} else if tx.Kind == "crypto_earn_interest_paid" ||
						tx.Kind == "crypto_earn_extra_interest_paid" ||
						tx.Kind == "mco_stake_reward" ||
						tx.Kind == "staking_reward" {
						cdc.TXsByCategory["Interests"] = append(cdc.TXsByCategory["Interests"], t)
					} else if tx.Kind == "referral_bonus" {
						cdc.TXsByCategory["Referrals"] = append(cdc.TXsByCategory["Referrals"], t)
					} else {
						cdc.TXsByCategory["Deposits"] = append(cdc.TXsByCategory["Deposits"], t)
					}
				} else if tx.Kind == "crypto_payment" ||
					tx.Kind == "crypto_withdrawal" ||
					tx.Kind == "card_cashback_reverted" ||
					tx.Kind == "transfer_cashback_reverted" ||
					tx.Kind == "reimbursement_reverted" ||
					tx.Kind == "crypto_to_exchange_transfer" ||
					tx.Kind == "supercharger_deposit" ||
					tx.Kind == "crypto_viban_exchange" {
					t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " " + tx.Kind + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					if tx.Kind == "crypto_withdrawal" &&
						tx.Description == "Withdraw BTC" {
						fee := decimal.New(3, -4) // 0.0003, is it always the case ? I have only one occurence
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: fee})
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg().Sub(fee)})
					} else {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
					}
					if tx.Kind == "crypto_payment" ||
						tx.Kind == "crypto_viban_exchange" {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.NativeCurrency, Amount: tx.NativeAmount.Neg()})
						cdc.TXsByCategory["CashOut"] = append(cdc.TXsByCategory["CashOut"], t)
					} else if tx.Kind == "card_cashback_reverted" ||
						tx.Kind == "transfer_cashback_reverted" ||
						tx.Kind == "reimbursement_reverted" {
						cdc.TXsByCategory["CommercialRebates"] = append(cdc.TXsByCategory["CommercialRebates"], t)
					} else {
						if is, desc := cat.IsTxGift(tx.ID); is {
							t.Note += " gift " + desc
							cdc.TXsByCategory["Gifts"] = append(cdc.TXsByCategory["Gifts"], t)
						} else {
							cdc.TXsByCategory["Withdrawals"] = append(cdc.TXsByCategory["Withdrawals"], t)
						}
					}
				} else if tx.Kind == "crypto_transfer" {
					t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " " + tx.Kind + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					if tx.Amount.IsNegative() {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
						if is, desc, val, curr := cat.IsTxCashOut(tx.ID); is {
							t.Note += " " + desc
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
							cdc.TXsByCategory["CashOut"] = append(cdc.TXsByCategory["CashOut"], t)
						} else {
							cdc.TXsByCategory["Gifts"] = append(cdc.TXsByCategory["Gifts"], t)
						}
					} else {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
						if is, desc, val, curr := cat.IsTxCashIn(tx.ID); is {
							t.Note += " " + desc
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
							cdc.TXsByCategory["CashIn"] = append(cdc.TXsByCategory["CashIn"], t)
						} else {
							cdc.TXsByCategory["Gifts"] = append(cdc.TXsByCategory["Gifts"], t)
						}
					}
				} else if tx.Kind == "crypto_earn_program_created" ||
					tx.Kind == "crypto_earn_program_withdrawn" ||
					tx.Kind == "lockup_lock" ||
					tx.Kind == "lockup_upgrade" ||
					tx.Kind == "lockup_swap_rebate" ||
					tx.Kind == "dynamic_coin_swap_bonus_exchange_deposit" ||
					tx.Kind == "dynamic_coin_swap_credited" ||
					tx.Kind == "dynamic_coin_swap_debited" ||
					tx.Kind == "viban_withdrawal" {
					// Do nothing
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Kind, tx, alreadyAsked)
				}
			}
		}
	}
	cdc.Sources["CdC App"] = source.Source{
		Crypto:        true,
		AccountNumber: account,
		OpeningDate:   firstTimeUsed,
		ClosingDate:   lastTimeUsed,
		LegalName:     "MCO Malta DAX Limited",
		Address:       "Level 7, Spinola Park, Triq Mikiel Ang Borg,\nSt Julian's SPK 1000,\nMalte",
		URL:           "https://crypto.com/app",
	}
	if hasCashback {
		switchGBLT := time.Date(2020, 12, 27, 3, 0, 0, 0, time.UTC)
		if firstTimeCashback.Before(switchGBLT) {
			if lastTimeCashback.Before(switchGBLT) {
				cdc.Sources["CdC MCO Card GB"] = source.Source{
					Crypto:        false,
					AccountNumber: "votre IBAN GBxxxxx",
					OpeningDate:   firstTimeCashback,
					ClosingDate:   lastTimeCashback,
					LegalName:     "MCO Malta DAX Limited The Currency Cloud)",
					Address:       "12 Steward Street, The Steward Building, London, E1 6FQ, Royaume-Uni",
					URL:           "https://crypto.com/cards",
				}
			} else {
				cdc.Sources["CdC MCO Card GB"] = source.Source{
					Crypto:        false,
					AccountNumber: "votre IBAN GBxxxxx",
					OpeningDate:   firstTimeCashback,
					ClosingDate:   switchGBLT,
					LegalName:     "MCO Malta DAX Limited The Currency Cloud)",
					Address:       "12 Steward Street, The Steward Building, London, E1 6FQ, Royaume-Uni",
					URL:           "https://crypto.com/cards",
				}
				cdc.Sources["CdC MCO Card LT"] = source.Source{
					Crypto:        false,
					AccountNumber: "votre IBAN LTxxxxx",
					OpeningDate:   switchGBLT,
					ClosingDate:   lastTimeCashback,
					LegalName:     "MCO Malta DAX Limited (Transactive Systems UAB)",
					Address:       "Jogailos St 9, Vilnius, 01103, Lithuania",
					URL:           "https://crypto.com/cards",
				}
			}
		} else {
			cdc.Sources["CdC MCO Card LT"] = source.Source{
				Crypto:        false,
				AccountNumber: "votre IBAN LTxxxxx",
				OpeningDate:   firstTimeCashback,
				ClosingDate:   lastTimeCashback,
				LegalName:     "MCO Malta DAX Limited (Transactive Systems UAB)",
				Address:       "Jogailos St 9, Vilnius, 01103, Lithuania",
				URL:           "https://crypto.com/cards",
			}
		}
	}
	return
}
