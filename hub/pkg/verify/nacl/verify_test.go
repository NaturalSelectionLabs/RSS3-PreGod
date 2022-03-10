package nacl_test

import (
	"encoding/base64"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/nacl"
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
			// Demo data from RSS3 v0.3.1 [0x1469Bd38ef9cF15917d052D6B10267ea041Ba091](https://prenode.rss3.dev/0x1469Bd38ef9cF15917d052D6B10267ea041Ba091)
			args: args{
				msg:    `[["agent_id","AEop+wb5O1eTBe5Q2DQ4ZAD0McOgggeWm6P4ZmNovCg="],["assets",[["list_auto","0x1469Bd38ef9cF15917d052D6B10267ea041Ba091-list-assets.auto-0"]]],["backlinks",[]],["date_created","2021-09-03T13:58:39.005Z"],["date_updated","2022-03-10T07:20:47.485Z"],["id","0x1469Bd38ef9cF15917d052D6B10267ea041Ba091"],["items",[["list_auto","0x1469Bd38ef9cF15917d052D6B10267ea041Ba091-list-items.auto-0"]]],["links",[["0",[["id","following"],["list","0x1469Bd38ef9cF15917d052D6B10267ea041Ba091-list-links.following-0"]]]]],["profile",[["accounts",[["0",[["id","Twitter-lc825"],["tags",[["0","pass:order:0"]]]]]]],["avatar",[["0","https://rss3.mypinata.cloud/ipfs/Qmakn9hSmDqqB7o5eN2j2GP6UrPUhFk6SzpVD9vP2EttiG"]]],["bio","Lemon!!!!"],["name","Lemon"]]],["version","rss3.io/version/v0.3.1"]]`, //nolint:lll // long str needed
				sig:    "VNrLH+SpbkjO1stOuH9u6mpk5Z57T42LZJWu+NXo4D++0SsdYdfdMT3XwMQQBq+9zUyYnaWta82lnUTXshQbCg==",
				pubkey: "AEop+wb5O1eTBe5Q2DQ4ZAD0McOgggeWm6P4ZmNovCg=",
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
