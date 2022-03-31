package verify_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify"
)

func TestSignature(t *testing.T) {
	type args struct {
		jsonBytes   []byte
		address     string
		instanceUri string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				jsonBytes:   []byte(`{}`),
				address:     "",
				instanceUri: "",
			},
			want: false,
		},
		{
			name: "random-1",
			args: args{
				jsonBytes:   []byte(`{"version":"v0.4.0","date_created":"2021-08-17T14:36:00.676Z","date_updated":"2022-02-10T22:50:53.132Z","agents":[{"pubkey":"rrqJ2xn7oUd4wGW8VbsZk9XeacYMap4/jprIA5b35ns=","signature":"PObUwUA+BEStJZJoY4xBsOkQujsRAZ4yULZIu0orDHCID2ezI5/eD8EskIK+RFNvSCp9tKTSYqurEFa2egW6Dg==","authorization":"","app":"Revery","date_expired":"2023-02-10T22:50:53.132Z"}],"profile":{"name":"DIYgod","avatars":["ipfs://QmT1zZNHvXxdTzHesfEdvFjMPvw536Ltbup7B4ijGVib7t"],"bio":"Cofounder of RSS3.","attachments":[{"type":"websites","content":"https://rss3.io\nhttps://diygod.me","mime_type":"text/uri-list"},{"type":"banner","content":"ipfs://QmT1zZNHvXxdTzHesfEdvFjMPvw536Ltbup7B4ijGVib7t","mime_type":"image/jpeg"}],"accounts":[{"identifier":"rss3://account:0x8768515270aA67C624d3EA3B98CA464672C50183@ethereum","signature":"0x4828da56a162b9504dca6009864a90ed0ca3e56256d8458af451874ad7dd9cb26be4f399a56a8b69a881297ba6b6434a7f2f4a4f3557890d1efa8490769187be1b"},{"identifier":"rss3://account:DIYgod@twitter"}]},"links":{"identifiers":[{"type":"following","identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following/1","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following"}],"identifier_back":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/backlink"},"items":{"notes":{"identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note/0","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note","filters":{"blocklist":["Twitter"]}},"assets":{"identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset/0","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset","filters":{"allowlist":["Polygon"]}}},"identifier":"rss3://account:0x8408AEd5907b21211A5f0915691F7DC1d9237328@ethereum","signature":"0x9dde74da5a55d77528b9c0f5d197d2609e643ca451d5525e3f921d4ab47d59596cbd2e311df514c804ee37c7a1fb275a930cf05fcb9169f81e00ac29c6025e491b"}`), // nolint:lll // need to be long
				address:     "0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944",
				instanceUri: "rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum",
			},
			want: false, // TODO ethers.VerifyMessage not OK
		},
		{
			name: "random-2-without-accounts",
			args: args{
				jsonBytes:   []byte(`{"version":"v0.4.0","date_created":"2021-08-17T14:36:00.676Z","date_updated":"2022-02-10T22:50:53.132Z","agents":[{"pubkey":"rrqJ2xn7oUd4wGW8VbsZk9XeacYMap4/jprIA5b35ns=","signature":"PObUwUA+BEStJZJoY4xBsOkQujsRAZ4yULZIu0orDHCID2ezI5/eD8EskIK+RFNvSCp9tKTSYqurEFa2egW6Dg==","authorization":"","app":"Revery","date_expired":"2023-02-10T22:50:53.132Z"}],"profile":{"name":"DIYgod","avatars":["ipfs://QmT1zZNHvXxdTzHesfEdvFjMPvw536Ltbup7B4ijGVib7t"],"bio":"Cofounder of RSS3.","attachments":[{"type":"websites","content":"https://rss3.io\nhttps://diygod.me","mime_type":"text/uri-list"},{"type":"banner","content":"ipfs://QmT1zZNHvXxdTzHesfEdvFjMPvw536Ltbup7B4ijGVib7t","mime_type":"image/jpeg"}]},"links":{"identifiers":[{"type":"following","identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following/1","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following"}],"identifier_back":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/backlink"},"items":{"notes":{"identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note/0","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note","filters":{"blocklist":["Twitter"]}},"assets":{"identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset/0","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset","filters":{"allowlist":["Polygon"]}}},"identifier":"rss3://account:0xA84836981A2C13f6078FF8706891a11264160117@ethereum","signature":"0x758b7433601e005511c2419795fc8ee87bb2020230129622b1c6d189b702e089297400240b06436f344696abe8ce6dba9d75142af5db478fc10605f6e6298d0c1b"}`), // nolint:lll // need to be long
				address:     "0xA84836981A2C13f6078FF8706891a11264160117",
				instanceUri: "rss3://account:0xA84836981A2C13f6078FF8706891a11264160117@ethereum",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := verify.Signature(tt.args.jsonBytes, tt.args.address, tt.args.instanceUri); got != tt.want {
				t.Errorf("Signature() = %v, want %v", got, tt.want)
			}
		})
	}
}
