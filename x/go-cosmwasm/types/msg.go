package types

import (
	"encoding/json"
)

//------- Results / Msgs -------------

// HandleResult is the raw response from the handle call.
// This is mirrors Rust's ContractResult<HandleResponse>.
type HandleResult struct {
	Ok  *HandleResponse `json:"ok,omitempty"`
	Err string          `json:"error,omitempty"`
}

// HandleResponse defines the return value on a successful handle
type HandleResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:"messages"`
	// base64-encoded bytes to return as ABCI.Data field
	Data []byte `json:"data"`
	// attributes for a log event to return over abci interface
	Attributes []EventAttribute `json:"attributes"`
}

// InitResult is the raw response from the handle call.
// This is mirrors Rust's ContractResult<InitResponse>.
type InitResult struct {
	Ok  *InitResponse `json:"ok,omitempty"`
	Err string        `json:"error,omitempty"`
}

// InitResponse defines the return value on a successful handle
type InitResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:"messages"`
	// attributes for a log event to return over abci interface
	Attributes []EventAttribute `json:"attributes"`
}

// MigrateResult is the raw response from the migrate call.
// This is mirrors Rust's ContractResult<MigrateResponse>.
type MigrateResult struct {
	Ok  *MigrateResponse `json:"ok,omitempty"`
	Err string           `json:"error,omitempty"`
}

// MigrateResponse defines the return value on a successful handle
type MigrateResponse struct {
	// Messages comes directly from the contract and is it's request for action
	Messages []CosmosMsg `json:"messages"`
	// base64-encoded bytes to return as ABCI.Data field
	Data []byte `json:"data"`
	// attributes for a log event to return over abci interface
	Attributes []EventAttribute `json:"attributes"`
}

// EventAttribute
type EventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CosmosMsg is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type CosmosMsg struct {
	Bank    *BankMsg        `json:"bank,omitempty"`
	Custom  json.RawMessage `json:"custom,omitempty"`
	Staking *StakingMsg     `json:"staking,omitempty"`
	Wasm    *WasmMsg        `json:"wasm,omitempty"`
	NFT		*NftMsg			`json:"nft,omitempty"`
}

type BankMsg struct {
	Send *SendMsg `json:"send,omitempty"`
}

// SendMsg contains instructions for a Cosmos-SDK/SendMsg
// It has a fixed interface here and should be converted into the proper SDK format before dispatching
type SendMsg struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      Coins  `json:"amount"`
}

type StakingMsg struct {
	Delegate   *DelegateMsg   `json:"delegate,omitempty"`
	Undelegate *UndelegateMsg `json:"undelegate,omitempty"`
	Redelegate *RedelegateMsg `json:"redelegate,omitempty"`
	Withdraw   *WithdrawMsg   `json:"withdraw,omitempty"`
}

type DelegateMsg struct {
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

type UndelegateMsg struct {
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

type RedelegateMsg struct {
	SrcValidator string `json:"src_validator"`
	DstValidator string `json:"dst_validator"`
	Amount       Coin   `json:"amount"`
}

type WithdrawMsg struct {
	Validator string `json:"validator"`
	// this is optional
	Recipient string `json:"recipient,omitempty"`
}

type WasmMsg struct {
	Execute     *ExecuteMsg     `json:"execute,omitempty"`
	Instantiate *InstantiateMsg `json:"instantiate,omitempty"`
}

// ExecuteMsg is used to call another defined contract on this chain.
// The calling contract requires the callee to be defined beforehand,
// and the address should have been defined in initialization.
// And we assume the developer tested the ABIs and coded them together.
//
// Since a contract is immutable once it is deployed, we don't need to transform this.
// If it was properly coded and worked once, it will continue to work throughout upgrades.
type ExecuteMsg struct {
	// ContractAddr is the sdk.AccAddress of the contract, which uniquely defines
	// the contract ID and instance ID. The sdk module should maintain a reverse lookup table.
	ContractAddr string `json:"contract_addr"`
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte `json:"msg"`
	// Send is an optional amount of coins this contract sends to the called contract
	Send Coins `json:"send"`
}

type InstantiateMsg struct {
	// CodeID is the reference to the wasm byte code as used by the Cosmos-SDK
	CodeID uint64 `json:"code_id"`
	// Msg is assumed to be a json-encoded message, which will be passed directly
	// as `userMsg` when calling `Handle` on the above-defined contract
	Msg []byte `json:"msg"`
	// Send is an optional amount of coins this contract sends to the called contract
	Send Coins `json:"send"`
}

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
}
