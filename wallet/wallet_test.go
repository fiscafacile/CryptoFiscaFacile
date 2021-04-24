package wallet

import (
	"testing"
	"time"
)

func TestWallet_TXsSortByDate(t *testing.T) {
	type args struct {
		chrono bool
	}
	tests := []struct {
		name     string
		args     args
		startTXs TXs
		wantTXs  TXs
	}{
		{
			name: "TXs Sort By Date order Chrono",
			args: args{
				chrono: true,
			},
			startTXs: TXs{
				TX{Timestamp: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
				TX{Timestamp: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantTXs: TXs{
				TX{Timestamp: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
				TX{Timestamp: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "TXs Sort By Date order anti Chrono",
			args: args{
				chrono: false,
			},
			startTXs: TXs{
				TX{Timestamp: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
				TX{Timestamp: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantTXs: TXs{
				TX{Timestamp: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
				TX{Timestamp: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.startTXs.SortByDate(tt.args.chrono)
			for i := range tt.startTXs {
				if !tt.startTXs[i].Timestamp.Equal(tt.wantTXs[i].Timestamp) {
					t.Errorf("TXs.SortByDate() startTXs = %v, wantTXs %v", tt.startTXs, tt.wantTXs)
				}
			}
		})
	}
}
