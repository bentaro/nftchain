package types

type NftMsg struct {
	Transfer   *TransferMsg   `json:"transfer,omitempty"`
	Mint *MintMsg `json:"mint,omitempty"`
	Burn *BurnMsg `json:"burn,omitempty"`
}

type TransferMsg struct {
	Sender string `json:"sender"`
	Recipient string `json:"recipient"`
	Denom string `json:"denom"`
	ID string `json:"id"`
}

type MintMsg struct {
	Sender string `json:"sender"`
	Recipient string `json:"recipient"`
	Denom string `json:"denom"`
	ID string `json:"id"`
	TokenURI string `json:"token_uri"`
}

type BurnMsg struct {
	Sender string `json:"sender"`
	Denom string `json:"denom"`
	ID string `json:"id"`
}