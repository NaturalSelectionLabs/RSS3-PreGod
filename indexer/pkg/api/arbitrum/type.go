package arbitrum

import "fmt"

type NFTItem struct {
	TokenAddress string `json:"token_address"`
	TokenId      string `json:"token_id"`
	ContractType string `json:"contract_type"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	TokenURI     string `json:"token_uri"`
	MetaData     string `json:"metadata"`
	Valid        bool   `json:"valid"`
}

func (i NFTItem) GetUid() string {
	return fmt.Sprintf("%s.%s", i.TokenAddress, i.TokenId)
}

func (i NFTItem) String() string {
	return fmt.Sprintf(`TokenAddress: %s, TokenId: %s, Type: %s, Name: %s, Symbol: %s, TokenURI: %s`,
		i.TokenAddress, i.TokenId, i.ContractType, i.Name, i.Symbol, i.TokenURI)
}

type NFTTransferItem struct {
	TokenAddress string
	Name         string
	Symbol       string
	TokenId      string
	From         string
	To           string
	TimeStamp    string
	Hash         string
}

func (i NFTTransferItem) GetUid() string {
	return fmt.Sprintf("%s.%s", i.TokenAddress, i.TokenId)
}

func (i NFTTransferItem) EqualsToToken(nft NFTItem) bool {
	return i.TokenAddress == nft.TokenAddress && i.TokenId == nft.TokenId
}
