package zksync

import "fmt"

type Token struct {
	Id       int64
	Address  string
	Symbol   string
	Decimals int64
	Kind     string
	IsNFT    bool
}

func (t Token) String() string {
	return fmt.Sprintf(`Id: %d, Address: %s, Symbol: %s, Decimals: %d, Kind: %s, IsNFT: %v`,
		t.Id, t.Address, t.Symbol, t.Decimals, t.Kind, t.IsNFT)
}

type Op struct {
	To        string
	Fee       string
	From      string
	Type      string
	Nonce     int64
	TokenId   int64
	Amount    string
	AccountId int64
}

func (o Op) String() string {
	return fmt.Sprintf(`From: %s, To: %s, Type: %s, TokenId: %d, Amount: %s`,
		o.From, o.To, o.Type, o.TokenId, o.Amount)
}

type ZKTransaction struct {
	TxHash      string
	BlockNumber int64
	Op          Op
	Success     bool
	FailReason  string
	CreatedAt   string
	BatchId     string
}
