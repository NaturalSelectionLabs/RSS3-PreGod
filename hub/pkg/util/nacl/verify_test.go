package nacl_test

import (
	"encoding/base64"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/util/nacl"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	type args struct {
		msg    string
		sig    string
		pubkey string
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
				msg:    `'[["_id",[["$oid","61c54ec1a3bae407078a1a68"]]],["content",[["_pass",[["assets",[["0",[["id","EVM+-0x000000000A38444e0a6E37d3b630d7e855a7cb13-Polygon.NFT-0xb4d21caf1cc3dadec5eecf753f5fc23094ddfb65.736"]]],["1",[["id","EVM+-0x000000000A38444e0a6E37d3b630d7e855a7cb13-Polygon.NFT-0xb766abf815d9e60a79ed96eac1093fbe46f31655.95314"],["order",1]]],["2",[["hide",true],["id","EVM+-0x000000000A38444e0a6E37d3b630d7e855a7cb13-Polygon.NFT-0xdf9293c820d5ac290b97fa62b449a9c790014296.1473"]]],["3",[["hide",true],["id","EVM+-0x000000000A38444e0a6E37d3b630d7e855a7cb13-Polygon.NFT-0x3c7174c46f6aeca7f485665a8ccfa29047f880e8.29"]]],["4",[["hide",true],["id","EVM+-0x000000000A38444e0a6E37d3b630d7e855a7cb13-Polygon.NFT-0x1fdf97e5bee48893eef28116973ca81166e4ec02.884"]]]]]]],["agent_id","NItLuV6qn0hybnExRU4Q42/ISGqYMpI2VHmbRcPdka0="],["assets",[["list_auto","0x000000000A38444e0a6E37d3b630d7e855a7cb13-list-assets.auto-0"]]],["backlinks",[]],["date_created","2021-12-20T15:56:03.737Z"],["date_updated","2021-12-25T07:57:34.114Z"],["id","0x000000000A38444e0a6E37d3b630d7e855a7cb13"],["items",[["list_auto","0x000000000A38444e0a6E37d3b630d7e855a7cb13-list-items.auto-0"]]],["profile",[["avatar",[["0","https://infura-ipfs.io/ipfs/QmcK8FSTtLQVydLEDKLv1hEacLxZgi7j2i4mkQQMyKxv6k"]]],["bio","Just a lazy man =U=<SITE#watertim.dev/>"],["name","WaterTim"]]],["version","rss3.io/version/v0.3.1"]]],["path","0x000000000A38444e0a6E37d3b630d7e855a7cb13"]]'`, //nolint:lll // long str needed
				sig:    "Kz47gL88TIJPnHcNBld5dLe7P005AFPnN/haqSPPyufb6/ZXdiQNvZVHI0ImpWUhhMPRnh6tGIZnpX45wcMCBQ==",
				pubkey: "MQ4rRFwbMa8BNsf8r+r/mCLT1qiZotWaayq//iMdbDc=",
			},
			want: true,
		},
		{
			name: "invalid",
			args: args{
				msg:    "invalid",
				sig:    base64.StdEncoding.EncodeToString([]byte("random-sig")),
				pubkey: base64.StdEncoding.EncodeToString([]byte("random-pubkey")),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			msg := []byte(tt.args.msg)

			sig, err := base64.StdEncoding.DecodeString(tt.args.sig)
			if err != nil && tt.want == true {
				t.Errorf("[%v] Verify() error = %v", tt.name, err)

				return
			}

			pubkey, err := base64.StdEncoding.DecodeString(tt.args.pubkey)
			if err != nil && tt.want == true {
				t.Errorf("[%v] Verify() error = %v", tt.name, err)

				return
			}

			got, err := nacl.Verify(msg, sig, pubkey)
			if err != nil && tt.want == true {
				t.Errorf("[%v] Verify() error = %v", tt.name, err)

				return
			}

			if got != tt.want {
				t.Errorf("[%v] Verify() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
