package cryptocom

import (
	"testing"
)

func Test_API_Exchange_sign(t *testing.T) {
	type args struct {
		body map[string]interface{}
	}
	tests := []struct {
		name      string
		args      args
		apiKey    string
		apiSecret string
		wantSig   string
	}{
		{
			name: "Test Signature on Official doc example",
			args: args{
				body: map[string]interface{}{
					"id":     11,
					"method": "private/get-order-detail",
					"params": map[string]interface{}{
						"order_id": "337843775021233500",
					},
					"nonce": 1619956517732,
				},
			},
			apiKey:    "API_KEY",
			apiSecret: "SECRET_KEY",
			wantSig:   "8c17b4cfbb7073a5453e348ecb1b20a1e709af910a7d1fe4b4569bcf29736e58",
		},
	}
	for _, tt := range tests {
		cdc := New()
		cdc.NewExchangeAPI(tt.apiKey, tt.apiSecret)
		t.Run(tt.name, func(t *testing.T) {
			cdc.apiEx.sign(tt.args.body)
			if tt.args.body["sig"] != tt.wantSig {
				t.Errorf("TXs.SortByDate() body = %v, wantSig %v", tt.args.body, tt.wantSig)
			}
		})
	}
}
