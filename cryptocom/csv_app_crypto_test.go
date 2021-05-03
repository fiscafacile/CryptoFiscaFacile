package cryptocom

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
			name:    "ParseCSVAppCrypto referral_card_cashback",
			csv:     "2020-12-31 15:43:19,Card Cashback,CRO,26.96195063,,,EUR,1.27,1.5174771149,referral_card_cashback",
			wantErr: false,
		},
		{
			name:    "ParseCSVAppCrypto exchange_to_crypto_transfer",
			csv:     "2020-11-05 11:45:03,Transfer: Exchange -> App wallet,CRO,50.47244468,,,EUR,3.35,4.0027939645,exchange_to_crypto_transfer",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc := New()
			err := cdc.ParseCSVAppCrypto(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("CryptoCom.ParseCSVAppCrypto() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
