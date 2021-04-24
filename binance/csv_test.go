package binance

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
			name:    "ParseCSV Buy",
			csv:     "2020-04-13 09:22:50,Spot,Buy,ETH,0.75673607,\"\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Sell",
			csv:     "2020-04-13 09:22:50,Spot,Sell,BNB,-8.29000000,\"\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdraw",
			csv:     "2020-04-13 09:33:17,Spot,Withdraw,ETH,-0.75597933,Withdraw fee is included",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New()
			err := b.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Binance.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
