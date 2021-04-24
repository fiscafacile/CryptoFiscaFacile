package revolut

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
			name:    "ParseCSV CashIn",
			csv:     "Completed Date,Description,Paid Out (BTC),Paid In (BTC),Exchange Out, Exchange In, Balance (BTC), Category, Notes\n26 nov. 2019,Échanger to BTC FX Rate 1 ₿ = 6530.9001 €,,0.01,EUR 65.31,,0.057,Général,",
			wantErr: false,
		},
		{
			name:    "ParseCSV CashOut",
			csv:     "Completed Date,Description,Paid Out (BTC),Paid In (BTC),Exchange Out, Exchange In, Balance (BTC), Category, Notes\n15 févr. 2020,Échanger BTC to FX Rate 1 ₿ = 9297.4833 €,0.057,,BTC 0.057,,0.00,Général,",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revo := New()
			err := revo.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Revolut.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
