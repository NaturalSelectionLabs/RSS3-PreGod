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
		// TODO fix it
		// {
		// 	name: "random-1",
		// 	args: args{
		// nolint: lll // 		jsonBytes:   []byte(`{"version":"v0.4.0","date_created":"2021-08-17T14:36:00.676Z","date_updated":"2022-02-10T22:50:53.132Z","signature":"0xd31375ff3614a97f32551df4b566fdf907e37c66ca8e3e50f3ce35620f3007dc6adbc9bca6aea21b3da06a0652020bfca40bd53e798ad804c57106f5cf8ddd2b1b","agents":[{"pubkey":"rrqJ2xn7oUd4wGW8VbsZk9XeacYMap4/jprIA5b35ns=","signature":"PObUwUA+BEStJZJoY4xBsOkQujsRAZ4yULZIu0orDHCID2ezI5/eD8EskIK+RFNvSCp9tKTSYqurEFa2egW6Dg==","authorization":"","app":"Revery","date_expired":"2023-02-10T22:50:53.132Z"}],"profile":{"name":"DIYgod","avatars":["ipfs://QmT1zZNHvXxdTzHesfEdvFjMPvw536Ltbup7B4ijGVib7t"],"bio":"Cofounder of RSS3.","attachments":[{"type":"websites","content":"https://rss3.io\\nhttps://diygod.me","mime_type":"text/uri-list"},{"type":"banner","content":"ipfs://QmT1zZNHvXxdTzHesfEdvFjMPvw536Ltbup7B4ijGVib7t","mime_type":"image/jpeg"}],"accounts":[{"identifier":"rss3://account:0x8768515270aA67C624d3EA3B98CA464672C50183@ethereum","signature":"0x4828da56a162b9504dca6009864a90ed0ca3e56256d8458af451874ad7dd9cb26be4f399a56a8b69a881297ba6b6434a7f2f4a4f3557890d1efa8490769187be1b"},{"identifier":"rss3://account:DIYgod@twitter"}]},"links":{"identifiers":[{"type":"following","identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following/1","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following"}],"identifier_back":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/backlink"},"items":{"notes":{"identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note/0","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note","filters":{"blocklist":["Twitter"]}},"assets":{"identifier_custom":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset/0","identifier":"rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset","filters":{"allowlist":["Polygon"]}}},"identifier":"rss3://account:0xAFd9C673cD7D22498EcA59Ed19D7b44DFB1036fb@ethereum"}`), // nolint:lll
		// 		address:     "0xAFd9C673cD7D22498EcA59Ed19D7b44DFB1036fb",
		// 		instanceUri: "rss3://account:0xAFd9C673cD7D22498EcA59Ed19D7b44DFB1036fb@ethereum",
		// 	},
		// 	want: true,
		// },
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
