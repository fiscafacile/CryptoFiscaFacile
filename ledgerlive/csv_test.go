package ledgerlive

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
			csv:     "2019-12-02T09:51:31.000Z,BTC,IN,2.8949079,0.0001922,f978fc4c94e054fa473ac4099f584bd9c0c58ade49f1a0c5941fad3a180e54e6,Bitcoin,xpub6DCB4S5L5Mp4Whnd6waASfGnXDLZXqBGRTRZ45QHqXgreimLRiYyen6HYuCqxgZDwCAW5AE1DgTy5RKSsYdhFGcueUSNH9vbTvWZTmEkh2Z",
			wantErr: false,
		},
		{
			name:    "ParseCSV Withdrawal",
			csv:     "2020-05-01T09:28:20.000Z,ETH,OUT,0.000336789,0.000336789,0x83066f0381f81a11c8c2849312252bc2e23ce957ec46db90583699d89d4238da,Ethereum,xpub6DXuQW1FgeHbgaWqEH46Mk1v5E8sGyNRcJ8zx6945Mys9YRy7ZsqGPPhDhJCWM4rYAk6JR6PqJosRn9sJFyWBHWEoEPHES7eg9x7tkddNxs",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := New()
			err := ll.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("LedgerLive.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
