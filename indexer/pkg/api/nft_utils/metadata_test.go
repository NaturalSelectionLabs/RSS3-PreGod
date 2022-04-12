package nft_utils_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/nft_utils"
)

func TestParseNFTMetadata(t *testing.T) {
	t.Parallel()

	type args struct {
		metadata string
	}

	tests := []struct {
		name    string
		args    args
		want    nft_utils.Metadata
		wantErr bool
	}{
		{
			name: "Test ParseNFTMetadata",
			args: args{`
{
  "name": "Test NFT",
  "description": "Test NFT description",
  "external_link": "https://example.com",
  "image": "https://example.com/image.png",
  "animation_url": "https://example.com/animation.gif",
  "attributes": [
    {
      "trait_type": "Base", 
      "value": "Starfish"
    }, 
    {
      "trait_type": "Eyes", 
      "value": "Big"
    }
  ]
}`,
			},
			want: nft_utils.Metadata{
				Name:         "Test NFT",
				Description:  "Test NFT description",
				ExternalLink: "https://example.com",
				Preview:      "https://example.com/image.png",
				Object:       "https://example.com/animation.gif",
				Attributes:   `[{"trait_type":"Base","value":"Starfish"},{"trait_type":"Eyes","value":"Big"}]`,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := nft_utils.ParseNFTMetadata(tt.args.metadata)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNFTMetadata() error = %+v, wantErr %+v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && !metadataEqual(got, tt.want) {
				t.Errorf("ParseNFTMetadata() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func metadataEqual(m1, m2 nft_utils.Metadata) bool {
	return m1.Name == m2.Name &&
		m1.Description == m2.Description &&
		m1.ExternalLink == m2.ExternalLink &&
		m1.Preview == m2.Preview &&
		m1.Object == m2.Object &&
		m1.Attributes == m2.Attributes
}
