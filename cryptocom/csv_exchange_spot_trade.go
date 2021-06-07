package cryptocom

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvSpotTrade struct {
	txsByCategory wallet.TXsByCategory
}

type csvExSpotTradeTX struct {
	AccountType        string
	OrderID            string
	TradeID            string
	CreateTimeUTC      time.Time
	SymbolLeft         string
	SymbolRight        string
	Side               string
	LiquidityIndicator string
	TradedPrice        decimal.Decimal
	TradedQuantity     decimal.Decimal
	Fee                decimal.Decimal
	FeeCurrency        string
}

func (cdc *CryptoCom) ParseCSVExchangeSpotTrade(reader io.Reader) (err error) {
	const SOURCE = "Crypto.com Exchange Spot Trade CSV :"
	alreadyAsked := []string{}
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "account_type" {
				tx := csvExSpotTradeTX{}
				tx.AccountType = r[0]
				tx.OrderID = r[1]
				tx.TradeID = r[2]
				tx.CreateTimeUTC, err = time.Parse("2006-01-02 15:04:05.000", r[3])
				if err != nil {
					log.Println(SOURCE, "Error Parsing CreateTimeUTC", r[3])
				}
				symbol := strings.Split(r[4], "_")
				tx.SymbolLeft = symbol[0]
				tx.SymbolRight = symbol[1]
				tx.Side = r[5]
				tx.LiquidityIndicator = r[6]
				tx.TradedPrice, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println(SOURCE, "Error Parsing TradedPrice", r[7])
				}
				tx.TradedQuantity, err = decimal.NewFromString(r[8])
				if err != nil {
					log.Println(SOURCE, "Error Parsing TradedQuantity", r[8])
				}
				tx.Fee, err = decimal.NewFromString(r[9])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Fee", r[9])
				}
				tx.FeeCurrency = r[10]
				cdc.csvExSpotTradeTXs = append(cdc.csvExSpotTradeTXs, tx)
				// Fill txsByCategory
				t := wallet.TX{Timestamp: tx.CreateTimeUTC, ID: tx.OrderID + "-" + tx.TradeID, Note: SOURCE + " " + tx.Side + " " + tx.LiquidityIndicator}
				t.Items = make(map[string]wallet.Currencies)
				if !tx.Fee.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
				}
				if tx.Side == "BUY" {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.SymbolLeft, Amount: tx.TradedQuantity})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.SymbolRight, Amount: tx.TradedQuantity.Mul(tx.TradedPrice)})
					cdc.csvSpotTrade.txsByCategory["Exchanges"] = append(cdc.csvSpotTrade.txsByCategory["Exchanges"], t)
				} else if tx.Side == "SELL" {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.SymbolLeft, Amount: tx.TradedQuantity})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.SymbolRight, Amount: tx.TradedQuantity.Mul(tx.TradedPrice)})
					cdc.csvSpotTrade.txsByCategory["Exchanges"] = append(cdc.csvSpotTrade.txsByCategory["Exchanges"], t)
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Side, tx, alreadyAsked)
				}
			}
		}
	}
	return
}
