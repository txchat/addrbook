package executor

import "testing"

func Test_checkChainNameField(t *testing.T) {
	type args struct {
		fd string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				fd: "chain.BTY",
			},
			want: true,
		}, {
			name: "no prefix failed",
			args: args{
				fd: "BTY",
			},
			want: false,
		}, {
			name: "prefix _ failed",
			args: args{
				fd: "chain_BTY",
			},
			want: false,
		}, {
			name: "prefix failed",
			args: args{
				fd: "chain BTY",
			},
			want: false,
		}, {
			name: "trim failed",
			args: args{
				fd: "chain.BT Y",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkChainNameField(tt.args.fd); got != tt.want {
				t.Errorf("checkChainNameField() = %v, want %v", got, tt.want)
			}
		})
	}
}
