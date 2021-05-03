package kraken

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
			csv:     "LYXXX-XXXXXX-XXXL5X,QCXXXXX-PNXXXX-PBXXXX,2018-01-05 11:36:20,deposit,,currency,ZEUR,10.0000,0.0000,10.0000,\"\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdraw",
			csv:     "LYXXX-XXXXXX-XXXL5X,QCXXXXX-PNXXXX-PBXXXX,2018-01-09 22:01:31,withdrawal,,currency,XLTC,-0.4700300000,0.0010000000,0.0010000000,\"\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Trade",
			csv:     "LYXXX-XXXXXX-XXXL5X,QCXXXXX-PNXXXX-PBXXXX,2018-01-09 13:21:13,trade,,currency,ZEUR,-149.7389,0.3893,350.7110,\"\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Staking",
			csv:     "LYXXX-XXXXXX-XXXL5X,QCXXXXX-PNXXXX-PBXXXX,2021-03-06 01:08:57,staking,,currency,DOT.S,0.0085120600,0.0000000000,10.0085120600,\"\"",
			wantErr: false,
		},
		{
			name:    "ParseCSV Transfer",
			csv:     "LYXXX-XXXXXX-XXXL5X,QCXXXXX-PNXXXX-PBXXXX,2021-03-02 10:49:32,transfer,stakingfromspot,currency,DOT.S,10.0000000000,0.0000000000,10.0000000000,\"\"",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := New()
			err := kr.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Kraken.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
