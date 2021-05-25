package uphold

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
			csv:     "Date,Destination,Destination Amount,Destination Currency,Fee Amount,Fee Currency,Id,Origin,Origin Amount,Origin Currency,Status,Type\nFri Apr 09 2021 23:58:41 GMT+0000,uphold,21.375,BAT,,,f0fc27a7-a0af-4a1a-8d42-24795a27f8fe,uphold,21.375,BAT,completed,in",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdrawal",
			csv:     "Date,Destination,Destination Amount,Destination Currency,Fee Amount,Fee Currency,Id,Origin,Origin Amount,Origin Currency,Status,Type\nWed May 05 2021 11:48:48 GMT+0000,uphold,5,BAT,,,7271cd79-b02d-4f52-b6dd-e663dc7c4cf9,uphold,5,BAT,completed,out",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := New()
			err := uh.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Uphold.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
