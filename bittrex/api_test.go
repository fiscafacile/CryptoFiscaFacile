package bittrex

import (
	"testing"
)

func TestAPI_TXsSortByHeight(t *testing.T) {
	type args struct {
		asc bool
	}
	tests := []struct {
		name     string
		args     args
		startTXs apiTXs
		wantTXs  apiTXs
	}{
		{
			name: "TXs Sort By Date order asc",
			args: args{
				asc: true,
			},
			startTXs: apiTXs{
				apiTX{Status: apiTXStatus{BlockHeight: 200}},
				apiTX{Status: apiTXStatus{BlockHeight: 100}},
			},
			wantTXs: apiTXs{
				apiTX{Status: apiTXStatus{BlockHeight: 100}},
				apiTX{Status: apiTXStatus{BlockHeight: 200}},
			},
		},
		{
			name: "TXs Sort By Height order desc",
			args: args{
				asc: false,
			},
			startTXs: apiTXs{
				apiTX{Status: apiTXStatus{BlockHeight: 100}},
				apiTX{Status: apiTXStatus{BlockHeight: 200}},
			},
			wantTXs: apiTXs{
				apiTX{Status: apiTXStatus{BlockHeight: 200}},
				apiTX{Status: apiTXStatus{BlockHeight: 100}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.startTXs.SortByHeight(tt.args.asc)
			for i := range tt.startTXs {
				if tt.startTXs[i].Status.BlockHeight != tt.wantTXs[i].Status.BlockHeight {
					t.Errorf("TXs.SortByDate() startTXs = %v, wantTXs %v", tt.startTXs, tt.wantTXs)
				}
			}
		})
	}
}
