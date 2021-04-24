package bluewallet

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
			csv:     "55edf8f6f8a3142d3e514596a9d05e863971753b435c901e821ed8d5d30308f5,1550997759,DEPOSIT,0.1,221,0.0000067",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdrawal",
			csv:     "51d96704b2357fcc8f258f0b536a672d90eb2d5ca029f421204537414e12878f,1560840474,WITHDRAWAL,0.04396333,249,0.0002276",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := New()
			err := bw.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueWallet.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
