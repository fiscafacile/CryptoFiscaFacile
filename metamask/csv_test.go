package metamask

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
			csv:     "9987949,0xde1638154ed2966de4b040cfff7cae2e07172b4948cec434f3934c960259d9a3,1588438154,53.23037922157064,sUSD,DEPOSIT,0x0000000000000000000000000000000000000000,0x0306fbf726fa310857ef7560fb6d3d2db42c32ce,1361705,1099738,5000000000",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdrawal",
			csv:     "9998583,0x3f99d996c6e0d779cd9c191d3c5553987469f10b06cdadd4dab65d5c5860e541,1588580259,0.40819096434933044,sUSD,WITHDRAWAL,0x0306fbf726fa310857ef7560fb6d3d2db42c32ce,0x0000000000000000000000000000000000000000,1209366,977548,8000000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mm := New()
			err := mm.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("MetaMask.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
