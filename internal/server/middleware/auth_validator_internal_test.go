package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBadHash(t *testing.T) {
	type args struct {
		cfgKey     string
		incomeHash string
		reqBody    []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "wrong",
			args: args{
				cfgKey:     "sdfs",
				incomeHash: "wrongHash",
				reqBody:    []byte("ABCDF"),
			},
			want: true,
		},
		{
			name: "ok",
			args: args{
				cfgKey:     "sdfs",
				incomeHash: "4f21bcb3b22ef261c261d033af3a8ad1fe8651f7edcc31180ac86a45b9040ee3",
				reqBody:    []byte("ABCDF"),
			},
			want: false,
		},
	}
	for _, unit := range tests {
		tt := unit
		t.Run(tt.name, func(t *testing.T) {
			got := isBadHash(tt.args.cfgKey, tt.args.incomeHash, tt.args.reqBody)
			assert.Equalf(t, tt.want, got, "isBadHash(%v, %v, %v)", tt.args.cfgKey, tt.args.incomeHash, tt.args.reqBody)
		})
	}
}
