package localbitcoin

import (
	"strings"
	"testing"
)

func Test_CSVParseTradeExemple(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		wantErr bool
	}{
		{
			name:    "ParseTradeCSV CashIn",
			csv:     "id,created_at,buyer,seller,trade_type,btc_amount,btc_traded,fee_btc,btc_amount_less_fee,btc_final,fiat_amount,fiat_fee,fiat_per_btc,currency,exchange_rate,transaction_released_at,online_provider,reference\n53688525,2019-11-29 06:34:25+00:00,moquette31,honestrade,ONLINE_SELL,0.09999955,0.09999955,0.00,0.09999955,0.09999955,6011.30,0.00,60113.27,HKD,60113.27,2019-11-29 06:48:36+00:00,NATIONAL_BANK,L53688525BVYQBX",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := New()
			err := lb.ParseTradeCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalBitcoin.ParseTradeCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_CSVLBParseTransferExemple(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		wantErr bool
	}{
		{
			name:    "ParseTransferCSV CashIn",
			csv:     "TXID, Created, Received, Sent, TXtype, TXdesc, TXNotes\n,2019-11-29T07:29:43+00:00,,1.76351515,Send to address,32F5pyzpge5KEi3CNZV5z9kE8d9ciqkm8k,",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := New()
			err := lb.ParseTransferCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalBitcoin.ParseTransferCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
