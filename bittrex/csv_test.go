package bittrex

import (
	"strings"
	"testing"
)

func TestBittrex_ParseCSV(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		wantErr bool
	}{
		{
			name:    "ParseCSV Exchange LIMIT_BUY",
			csv:     "40a1adf3-e43d-4e34-8bc6-5d1f8g2c3z6d,BTC-GBYTE,12/22/2017 9:10:21 AM,LIMIT_BUY,0.03700010,6.32828762,0.00000000,0.00058503,0.23401959,0.03697992,False,,0.00000000,False,12/22/2017 9:10:23 AM,0,\n",
			wantErr: false,
		},
		{
			name:    "ParseCSV Exchange MARKET_BUY",
			csv:     "48545d51-8bc6-e43d-4e34-5d1f8g2c3z6d,BTC-MANA,1/15/2017 9:04:21 PM,MARKET_BUY,0.03705010,10.32828762,0.00000000,0.00088503,0.23401359,0.08697992,False,,0.00000000,False,1/15/2020 9:06:50 PM,0,\n",
			wantErr: false,
		},
		{
			name:    "ParseCSV Exchange LIMIT_SELL",
			csv:     "48545d51-8bc6-e43d-4e34-5d1f8g2c3z6d,BTC-MCO,1/9/2020 9:04:21 PM,LIMIT_SELL,0.03705010,10.32828762,0.00000000,0.00088503,0.23401359,0.08697992,False,,0.00000000,False,1/9/2020 9:06:50 PM,0,\n",
			wantErr: false,
		},
		{
			name:    "ParseCSV Exchange MARKET_SELL",
			csv:     "40a1adf3-e43d-4e34-8bc6-5d1f8g2c3z6d,BTC-GBYTE,12/22/2017 9:10:21 AM,MARKET_SELL,0.03700010,6.32828762,0.00000000,0.00058503,0.23401959,0.03697992,False,,0.00000000,False,12/22/2017 9:10:23 AM,0,\n",
			wantErr: false,
		},
		{
			name:    "ParseCSV Exchange CEILING_MARKET_BUY",
			csv:     "48545d51-8bc6-e43d-4e34-5d1f8g2c3z6d,BTC-MANA,1/15/2017 9:04:21 PM,CEILING_MARKET_BUY,0.03705010,10.32828762,0.00000000,0.00088503,0.23401359,0.08697992,False,,0.00000000,False,1/15/2020 9:06:50 PM,0,\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btrx := New()
			err := btrx.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("Bittrex.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
