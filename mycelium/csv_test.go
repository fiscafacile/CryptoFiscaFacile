package mycelium

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
			name:    "ParseCSV",
			csv:     "Compte 1,563c9ae8edca798b1d13eb4f167f4a8735385ad9dcec767a1bf0377e43bf3929,16Rp4mkpFY4rgSzX7VFFbmUuJSZymqz83c,2018-11-06T23:08Z,-0.00924295,Bitcoin,",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := New()
			err := mc.ParseCSV(strings.NewReader(tt.csv))
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCelium.ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
