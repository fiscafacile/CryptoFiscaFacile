package bitfinex

import (
	"strings"
	"testing"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

func Test_CSVParseExemple(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		wantErr bool
	}{
		{
			name:    "ParseCSV Withdrawal",
			csv:     "2809474008,Crypto Withdrawal fee on wallet exchange,BTC,-0.0004,0,01-05-20 18:41:57,exchange",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := New()
			err := bf.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Bitfinex.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_CSVParseExchangeExemple(t *testing.T) {
	tests := []struct {
		name          string
		csv           string
		wantErr       bool
		wantExchanges wallet.TXs
	}{
		{
			name:    "ParseCSV Exchange FeeToFrom",
			csv:     "787134003,Trading fees for 10.481 XMR (XMRBTC) @ 0.0159 on BFX (0.2%) on wallet exchange,XMR,-0.020962,10.460038,10-12-17 19:50:08,exchange\n787133984,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,XMR,10.481,10.481,10-12-17 19:50:08,exchange\n787133969,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,BTC,-0.16683656,1.33316344,10-12-17 19:50:08,exchange\n",
			wantErr: false,
			wantExchanges: wallet.TXs{
				{
					Timestamp: time.Date(2017, time.December, 10, 19, 50, 8, 0, time.UTC),
					Items: map[string][]wallet.Currency{
						"From": {{Code: "BTC", Amount: 0.16683656}},
						"To":   {{Code: "XMR", Amount: 10.481}},
						"Fee":  {{Code: "XMR", Amount: 0.020962}},
					},
					Note: "Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange",
				},
			},
		},
		{
			name:    "ParseCSV Exchange FeeFromTo",
			csv:     "787134003,Trading fees for 10.481 XMR (XMRBTC) @ 0.0159 on BFX (0.2%) on wallet exchange,XMR,-0.020962,10.460038,10-12-17 19:50:08,exchange\n787133969,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,BTC,-0.16683656,1.33316344,10-12-17 19:50:08,exchange\n787133984,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,XMR,10.481,10.481,10-12-17 19:50:08,exchange\n",
			wantErr: false,
			wantExchanges: wallet.TXs{
				{
					Timestamp: time.Date(2017, time.December, 10, 19, 50, 8, 0, time.UTC),
					Items: map[string][]wallet.Currency{
						"From": {{Code: "BTC", Amount: 0.16683656}},
						"To":   {{Code: "XMR", Amount: 10.481}},
						"Fee":  {{Code: "XMR", Amount: 0.020962}},
					},
					Note: "Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange",
				},
			},
		},
		{
			name:    "ParseCSV Exchange ToFeeFrom",
			csv:     "787133984,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,XMR,10.481,10.481,10-12-17 19:50:08,exchange\n787134003,Trading fees for 10.481 XMR (XMRBTC) @ 0.0159 on BFX (0.2%) on wallet exchange,XMR,-0.020962,10.460038,10-12-17 19:50:08,exchange\n787133969,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,BTC,-0.16683656,1.33316344,10-12-17 19:50:08,exchange\n",
			wantErr: false,
			wantExchanges: wallet.TXs{
				{
					Timestamp: time.Date(2017, time.December, 10, 19, 50, 8, 0, time.UTC),
					Items: map[string][]wallet.Currency{
						"From": {{Code: "BTC", Amount: 0.16683656}},
						"To":   {{Code: "XMR", Amount: 10.481}},
						"Fee":  {{Code: "XMR", Amount: 0.020962}},
					},
					Note: "Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange",
				},
			},
		},
		{
			name:    "ParseCSV Exchange ToFromFee",
			csv:     "787133984,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,XMR,10.481,10.481,10-12-17 19:50:08,exchange\n787133969,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,BTC,-0.16683656,1.33316344,10-12-17 19:50:08,exchange\n787134003,Trading fees for 10.481 XMR (XMRBTC) @ 0.0159 on BFX (0.2%) on wallet exchange,XMR,-0.020962,10.460038,10-12-17 19:50:08,exchange\n",
			wantErr: false,
			wantExchanges: wallet.TXs{
				{
					Timestamp: time.Date(2017, time.December, 10, 19, 50, 8, 0, time.UTC),
					Items: map[string][]wallet.Currency{
						"From": {{Code: "BTC", Amount: 0.16683656}},
						"To":   {{Code: "XMR", Amount: 10.481}},
						"Fee":  {{Code: "XMR", Amount: 0.020962}},
					},
					Note: "Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange",
				},
			},
		},
		{
			name:    "ParseCSV Exchange FromFeeTo",
			csv:     "787133969,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,BTC,-0.16683656,1.33316344,10-12-17 19:50:08,exchange\n787134003,Trading fees for 10.481 XMR (XMRBTC) @ 0.0159 on BFX (0.2%) on wallet exchange,XMR,-0.020962,10.460038,10-12-17 19:50:08,exchange\n787133984,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,XMR,10.481,10.481,10-12-17 19:50:08,exchange\n",
			wantErr: false,
			wantExchanges: wallet.TXs{
				{
					Timestamp: time.Date(2017, time.December, 10, 19, 50, 8, 0, time.UTC),
					Items: map[string][]wallet.Currency{
						"From": {{Code: "BTC", Amount: 0.16683656}},
						"To":   {{Code: "XMR", Amount: 10.481}},
						"Fee":  {{Code: "XMR", Amount: 0.020962}},
					},
					Note: "Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange",
				},
			},
		},
		{
			name:    "ParseCSV Exchange FromToFee",
			csv:     "787133969,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,BTC,-0.16683656,1.33316344,10-12-17 19:50:08,exchange\n787133984,Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange,XMR,10.481,10.481,10-12-17 19:50:08,exchange\n787134003,Trading fees for 10.481 XMR (XMRBTC) @ 0.0159 on BFX (0.2%) on wallet exchange,XMR,-0.020962,10.460038,10-12-17 19:50:08,exchange\n",
			wantErr: false,
			wantExchanges: wallet.TXs{
				{
					Timestamp: time.Date(2017, time.December, 10, 19, 50, 8, 0, time.UTC),
					Items: map[string][]wallet.Currency{
						"From": {{Code: "BTC", Amount: 0.16683656}},
						"To":   {{Code: "XMR", Amount: 10.481}},
						"Fee":  {{Code: "XMR", Amount: 0.020962}},
					},
					Note: "Exchange 10.481 XMR for BTC @ 0.015918 on wallet exchange",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := New()
			err := bf.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Bitfinex.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(bf.Accounts["Exchanges"]) != len(tt.wantExchanges) ||
				!bf.Accounts["Exchanges"][0].Timestamp.Equal(tt.wantExchanges[0].Timestamp) ||
				bf.Accounts["Exchanges"][0].Items["From"][0].Code != tt.wantExchanges[0].Items["From"][0].Code ||
				bf.Accounts["Exchanges"][0].Items["From"][0].Amount != tt.wantExchanges[0].Items["From"][0].Amount ||
				bf.Accounts["Exchanges"][0].Items["To"][0].Code != tt.wantExchanges[0].Items["To"][0].Code ||
				bf.Accounts["Exchanges"][0].Items["To"][0].Amount != tt.wantExchanges[0].Items["To"][0].Amount ||
				bf.Accounts["Exchanges"][0].Items["Fee"][0].Code != tt.wantExchanges[0].Items["Fee"][0].Code ||
				bf.Accounts["Exchanges"][0].Items["Fee"][0].Amount != tt.wantExchanges[0].Items["Fee"][0].Amount ||
				bf.Accounts["Exchanges"][0].Note != tt.wantExchanges[0].Note {
				t.Errorf("Bitfinex.ParseCSV() bf.Accounts[\"Exchanges\"] = %v, wantExchanges %v", bf.Accounts["Exchanges"], tt.wantExchanges)
			}
		})
	}
}
