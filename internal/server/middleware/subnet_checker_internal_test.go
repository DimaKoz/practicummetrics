package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTrusted(t *testing.T) {
	type args struct {
		trusted string
		realIP  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				realIP:  "192.168.1.1",
				trusted: "192.168.1.0/24",
			},
			want: true,
		},
		{
			name: "bad",
			args: args{
				realIP:  "192.178.1.1",
				trusted: "192.168.1.0/24",
			},
			want: false,
		},
	}
	for _, testCase := range tests {
		tt := testCase
		t.Run(tt.name, func(t *testing.T) {
			got, err := isTrusted(tt.args.trusted, tt.args.realIP)
			assert.NoError(t, err)
			if tt.want {
				assert.True(t, got)
			} else {
				assert.False(t, got)
			}
		})
	}
}
