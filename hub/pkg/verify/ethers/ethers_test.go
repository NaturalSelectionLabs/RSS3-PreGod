package ethers_test

import (
	"fmt"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/ethers"
)

func TestVerifyMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		msg  string
		sig  string
		addr string
	}

	tests := []struct {
		name string
		args args
		want bool
		err  error
	}{
		{
			name: "valid",
			args: args{
				msg:  fmt.Sprintf("Hi, RSS3. I'm your agent %s", "NItLuV6qn0hybnExRU4Q42/ISGqYMpI2VHmbRcPdka0="),
				sig:  "0xb5719893bf59b1ee95587293e7b32332259c2753d1acad775b8b162ac53d6199567b407d4d7a5f67ab2bbd9634af30e0d5458a6bfea63870ac9c6f231d0daaea1c",
				addr: "0x000000000A38444e0a6E37d3b630d7e855a7cb13",
			},
			want: true,
		},
		{
			name: "hello",
			args: args{
				msg:  "hello",
				sig:  "0x53edb561b0c1719e46e1e6bbbd3d82ff798762a66d0282a9adf47a114e32cbc600c248c247ee1f0fb3a6136a05f0b776db4ac82180442d3a80f3d67dde8290811c",
				addr: "0x829814B6E4dfeC4b703F2c6fDba28F1724094D11",
			},
			want: true,
		},
		{
			name: "invalid",
			args: args{
				msg:  "bad",
				sig:  "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				addr: "0x0000000000000000000000000000000000000000",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ethers.VerifyMessage(tt.args.msg, tt.args.sig, tt.args.addr)
			if err != nil {
				t.Errorf("[%v] VerifyMessage() error = %v", tt.name, err)

				return
			}
			if got != tt.want {
				t.Errorf("[%v] VerifyMessage() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
