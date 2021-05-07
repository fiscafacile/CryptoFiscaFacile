package binance

import (
	"testing"
)

func Test_API_sign(t *testing.T) {
	type args struct {
		body map[string]interface{}
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
			name: "Test Signature on Official doc example",
			args: args{
				body: map[string]interface{}{
					"endpoint": "wapi/v3/withdraw.html",
					"params": map[string]interface{}{
						"asset":      "ETH",
						"address":    "0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b",
						"addressTag": "1",
						"amount":     "1",
						"recvWindow": "5000",
						"name":       "addressName",
						"timestamp":  "1508396497000",
					},
				},
			},
			debug:     false,
			apiKey:    "vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A",
			apiSecret: "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j",
			wantSig:   "7bf8f11e7d683bb1218d1da3802aa4821ac7f74c1868848dfa8c92815552dd89",
		},
	}
	for _, tt := range tests {
		b := New()
		b.NewAPI(tt.apiKey, tt.apiSecret, tt.debug)
		t.Run(tt.name, func(t *testing.T) {
			b.api.sign(tt.args.body)
			if tt.args.body["sig"] != tt.wantSig {
				t.Errorf("TXs.SortByDate() body = %v, wantSig %v", tt.args.body, tt.wantSig)
			}
		})
	}
}
