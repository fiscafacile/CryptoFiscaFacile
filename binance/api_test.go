package binance

import (
	"testing"
)

func Test_API_sign(t *testing.T) {
	type args struct {
		queryParams map[string]string
	}
	tests := []struct {
		name      string
		args      args
		apiKey    string
		apiSecret string
		wantSig   string
		debug     bool
	}{
		{
			name: "Test Signature",
			args: args{
				queryParams: map[string]string{
					"asset":      "ETH",
					"address":    "0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b",
					"amount":     "1",
					"recvWindow": "5000",
					"name":       "test",
					"timestamp":  "1510903211000",
				},
			},
			debug:     false,
			apiSecret: "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j",
			wantSig:   "89294eb90c24e5b1736723bd6afb03388b95ec793a359f49dca208142645839a",
		},
	}
	for _, tt := range tests {
		b := New()
		b.NewAPI(tt.apiKey, tt.apiSecret, tt.debug)
		t.Run(tt.name, func(t *testing.T) {
			b.api.sign(tt.args.queryParams)
			if tt.args.queryParams["signature"] != tt.wantSig {
				t.Errorf("TXs.SortByDate() queryParams = %v, wantSig %v", tt.args.queryParams, tt.wantSig)
			}
		})
	}
}
