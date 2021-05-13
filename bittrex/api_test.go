package bittrex

import (
	"testing"
)

func TestBittrexAPI_sign(t *testing.T) {
	type args struct {
		apiKey            string
		secretKey         string
		timestamp         string
		ressource         string
		method            string
		hash              string
		queryParamEncoded string
	}
	tests := []struct {
		name     string
		args     args
		wantSign string
	}{
		{
			name: "sign deposit",
			args: args{
				apiKey:            "apiKey",
				secretKey:         "apiSecret",
				timestamp:         "1620319842000",
				ressource:         "deposits/closed",
				method:            "GET",
				hash:              "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
				queryParamEncoded: "?status=COMPLETED&pageSize=200",
			},
			wantSign: "d412349b312138676da649a5f9c14a286eeca3cc5d6d29d5ce6ec7cac60b50d0256f385423442034be67202cd0b9ae043ee8bc45c383d9dfdde58fadceded14a",
		},
		{
			name: "sign orders",
			args: args{
				apiKey:            "apiKey",
				secretKey:         "apiSecret",
				timestamp:         "1620319842000",
				ressource:         "orders/closed",
				method:            "GET",
				hash:              "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
				queryParamEncoded: "?pageSize=200",
			},
			wantSign: "53217db44dad750db36c16cc815d7aa23c586d52dfe79bc79c36bace29dd53b7595d803505c4e2518b110fd0d1bfa10c0fa6bf18448835063b8f0f84a4931304",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btrx := New()
			btrx.NewAPI(tt.args.apiKey, tt.args.secretKey, false)
			_, sig := btrx.api.sign(tt.args.timestamp, tt.args.ressource, tt.args.method, tt.args.hash, tt.args.queryParamEncoded)
			if sig != tt.wantSign {
				t.Errorf("Bittrex.sign() sig = %v, wantSign %v", sig, tt.wantSign)
			}
		})
	}
}
