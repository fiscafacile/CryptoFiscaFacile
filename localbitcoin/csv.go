package localbitcoin

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/utils"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTXTrade struct {
	ID                    string
	CreatedAt             time.Time
	Buyer                 string
	Seller                string
	TradeType             string
	Amount                decimal.Decimal
	Traded                decimal.Decimal
	FeeBTC                decimal.Decimal
	AmountLessFee         decimal.Decimal
	Final                 decimal.Decimal
	FiatAmount            decimal.Decimal
	FiatFee               decimal.Decimal
	FiatPerBTC            decimal.Decimal
	FiatCurrency          string
	ExchangeRate          decimal.Decimal
	TransactionReleasedAt time.Time
	OnlineProvider        string
	Reference             string
}

type CsvTXTransfer struct {
	ID       string
	Created  time.Time
	Received decimal.Decimal
	Sent     decimal.Decimal
	Type     string
	Desc     string
	Notes    string
}

func (lb *LocalBitcoin) ParseTradeCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Local Bitcoin CSV Trade :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		var curr string
		for _, r := range records {
			if r[0] == "id" {
				curr = strings.Split(r[5], "_")[0]
				curr = strings.ToUpper(curr)
			} else {
				tx := CsvTXTrade{}
				tx.ID = r[0]
				tx.CreatedAt, err = time.Parse("2006-01-02 15:04:05+00:00", r[1])
				if err != nil {
					log.Println(SOURCE, "Error Parsing CreatedAt : ", r[1])
				}
				tx.Buyer = r[2]
				tx.Seller = r[3]
				tx.TradeType = r[4]
				tx.Amount, err = decimal.NewFromString(r[5])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount : ", r[5])
				}
				tx.Traded, err = decimal.NewFromString(r[6])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Traded : ", r[6])
				}
				tx.FeeBTC, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println(SOURCE, "Error Parsing FeeBTC : ", r[7])
				}
				tx.AmountLessFee, err = decimal.NewFromString(r[8])
				if err != nil {
					log.Println(SOURCE, "Error Parsing AmountLessFee : ", r[8])
				}
				tx.Final, err = decimal.NewFromString(r[9])
				if err != nil {
					log.Println(SOURCE, "Error Parsing BTC_Final : ", r[9])
				}
				tx.FiatAmount, err = decimal.NewFromString(r[10])
				if err != nil {
					log.Println(SOURCE, "Error Parsing FiatAmount : ", r[10])
				}
				tx.FiatFee, err = decimal.NewFromString(r[11])
				if err != nil {
					log.Println(SOURCE, "Error Parsing FiatFee : ", r[11])
				}
				tx.FiatPerBTC, err = decimal.NewFromString(r[12])
				if err != nil {
					log.Println(SOURCE, "Error Parsing FiatPerBTC : ", r[12])
				}
				tx.FiatCurrency = r[13]
				tx.ExchangeRate, err = decimal.NewFromString(r[14])
				if err != nil {
					log.Println(SOURCE, "Error Parsing ExchangeRate : ", r[14])
				}
				tx.TransactionReleasedAt, err = time.Parse("2006-01-02 15:04:05+00:00", r[15])
				if err != nil {
					log.Println(SOURCE, "Error Parsing TransactionReleasedAt : ", r[15])
				}
				tx.OnlineProvider = r[16]
				tx.Reference = r[17]
				lb.CsvTXsTrade = append(lb.CsvTXsTrade, tx)
				if tx.TransactionReleasedAt.Before(firstTimeUsed) {
					firstTimeUsed = tx.TransactionReleasedAt
				}
				if tx.TransactionReleasedAt.After(lastTimeUsed) {
					lastTimeUsed = tx.TransactionReleasedAt
				}
				// Fill TXsByCategory
				if tx.TradeType == "ONLINE_SELL" {
					t := wallet.TX{Timestamp: tx.TransactionReleasedAt, ID: tx.ID, Note: SOURCE + " " + tx.Seller + " " + tx.Buyer + " " + tx.TradeType + " " + tx.OnlineProvider + " " + tx.Reference}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: tx.Amount})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.FiatCurrency, Amount: tx.FiatAmount})
					if !tx.FiatFee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FiatCurrency, Amount: tx.FiatFee})
					}
					if !tx.FeeBTC.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: curr, Amount: tx.FeeBTC})
					}
					lb.TXsByCategory["Exchanges"] = append(lb.TXsByCategory["Exchanges"], t)
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.TradeType, tx, alreadyAsked)
				}
			}
		}
	}
	if _, ok := lb.Sources["Local Bitcoin"]; !ok {
		lb.Sources["Local Bitcoin"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "LocalBitcoins Oy",
			Address:       "Porkkalankatu 24\n00180 Helsinki\nFinland",
			URL:           "https://localbitcoins.com/fr",
		}
	}
	return
}

func (lb *LocalBitcoin) ParseTransferCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Local Bitcoin CSV Transfer :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		curr := "BTC"
		for _, r := range records {
			if r[0] != "TXID" {
				tx := CsvTXTransfer{}
				tx.Created, err = time.Parse("2006-01-02T15:04:05+00:00", r[1])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Created : ", r[1])
				}
				if r[0] != "" {
					tx.ID = r[0]
				} else {
					utils.GetUniqueID(SOURCE + tx.Created.String())
				}
				if r[2] != "" {
					tx.Received, err = decimal.NewFromString(r[2])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Received : ", r[2])
					}
				}
				if r[3] != "" {
					tx.Sent, err = decimal.NewFromString(r[3])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Sent : ", r[3])
					}
				}
				tx.Type = r[4]
				tx.Desc = r[5]
				tx.Notes = r[6]
				lb.CsvTXsTransfer = append(lb.CsvTXsTransfer, tx)
				if tx.Created.Before(firstTimeUsed) {
					firstTimeUsed = tx.Created
				}
				if tx.Created.After(lastTimeUsed) {
					lastTimeUsed = tx.Created
				}
				// Fill TXsByCategory
				if tx.Type == "Send to address" {
					t := wallet.TX{Timestamp: tx.Created, ID: tx.ID, Note: SOURCE + " " + tx.Type + " " + tx.Desc + " " + tx.Notes}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: tx.Sent})
					lb.TXsByCategory["Withdrawals"] = append(lb.TXsByCategory["Withdrawals"], t)
				} else if tx.Type == "Other" &&
					tx.Desc == "fee" {
					found := false
					for i, ex := range lb.TXsByCategory["Withdrawals"] {
						if ex.SimilarDate(2*time.Second, tx.Created) {
							found = true
							lb.TXsByCategory["Withdrawals"][i].Items["Fee"] = append(lb.TXsByCategory["Withdrawals"][i].Items["Fee"], wallet.Currency{Code: curr, Amount: tx.Sent})
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Created, ID: tx.ID, Note: SOURCE + " " + tx.Type + " " + tx.Desc + " " + tx.Notes}
						t.Items = make(map[string]wallet.Currencies)
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: curr, Amount: tx.Sent})
						lb.TXsByCategory["Withdrawals"] = append(lb.TXsByCategory["Withdrawals"], t)
					}
				} else if tx.Type == "Other" {
					// Do Nothing
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
				}
			}
		}
	}
	if _, ok := lb.Sources["Local Bitcoin"]; !ok {
		lb.Sources["Local Bitcoin"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "LocalBitcoins Oy",
			Address:       "Porkkalankatu 24\n00180 Helsinki\nFinland",
			URL:           "https://localbitcoins.com/fr",
		}
	}
	return
}
