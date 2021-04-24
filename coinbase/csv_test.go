package coinbase

import (
	"strings"
	"testing"
)

func Test_CSVParseExemple(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		wantErr bool
	}{
		{
			name:    "ParseCSV Deposit",
			csv:     "Timestamp,Transaction Type,Asset,Quantity Transacted,EUR Spot Price at Transaction,EUR Subtotal,EUR Total (inclusive of fees),EUR Fees,Notes\n2017-01-04T17:34:46Z,Receive,BTC,0.00979287,0.00,\"\",\"\",\"\",\"Received 0,00979287 BTC from an external account\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV CashOut",
			csv:     "Timestamp,Transaction Type,Asset,Quantity Transacted,EUR Spot Price at Transaction,EUR Subtotal,EUR Total (inclusive of fees),EUR Fees,Notes\n2017-08-10T13:26:40Z,Sell,BTC,0.00979287,2909.26,28.49,26.50,1.99,\"Sold 0,00979287 BTC for 26,50 € EUR\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV CashIn",
			csv:     "Timestamp,Transaction Type,Asset,Quantity Transacted,EUR Spot Price at Transaction,EUR Subtotal,EUR Total (inclusive of fees),EUR Fees,Notes\n2017-08-16T06:37:40Z,Buy,BTC,0.00728461,3431.89,25.00,26.49,1.49,\"Bought 0,00728461 BTC for 26,49 € EUR\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdrawals",
			csv:     "Timestamp,Transaction Type,Asset,Quantity Transacted,EUR Spot Price at Transaction,EUR Subtotal,EUR Total (inclusive of fees),EUR Fees,Notes\n2017-11-08T13:55:22Z,Send,BTC,0.15465874,6460.85,\"\",\"\",\"\",\"Sent 0,15465874 BTC to 1FBkm4BVF1bL164KD3sz4WXGwnkULZaG6X\"",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := New()
			err := cb.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Coinbase.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
