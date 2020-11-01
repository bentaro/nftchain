package cosmwasm

import (
	"encoding/json"
	"fmt"

	"github.com/CosmWasm/go-cosmwasm/api"
	"github.com/CosmWasm/go-cosmwasm/types"
)

// CodeID represents an ID for a given wasm code blob, must be generated from this library
type CodeID []byte

// WasmCode is an alias for raw bytes of the wasm compiled code
type WasmCode []byte

// KVStore is a reference to some sub-kvstore that is valid for one instance of a code
type KVStore = api.KVStore

// GoAPI is a reference to some "precompiles", go callbacks
type GoAPI = api.GoAPI

// Querier lets us make read-only queries on other modules
type Querier = types.Querier

// GasMeter is a read-only version of the sdk gas meter
type GasMeter = api.GasMeter

// Wasmer is the main entry point to this library.
// You should create an instance with it's own subdirectory to manage state inside,
// and call it for all cosmwasm code related actions.
type Wasmer struct {
	cache api.Cache
}

// NewWasmer creates an new binding, with the given dataDir where
// it can store raw wasm and the pre-compile cache.
// cacheSize sets the size of an optional in-memory LRU cache for prepared VMs.
// They allow popular contracts to be executed very rapidly (no loading overhead),
// but require ~32-64MB each in memory usage.
func NewWasmer(dataDir string, supportedFeatures string) (*Wasmer, error) {
	cache, err := api.InitCache(dataDir, supportedFeatures)
	if err != nil {
		return nil, err
	}
	return &Wasmer{cache: cache}, nil
}

// Cleanup should be called when no longer using this to free resources on the rust-side
func (w *Wasmer) Cleanup() {
	api.ReleaseCache(w.cache)
}

// Create will compile the wasm code, and store the resulting pre-compile
// as well as the original code. Both can be referenced later via CodeID
// This must be done one time for given code, after which it can be
// instatitated many times, and each instance called many times.
//
// For example, the code for all ERC-20 contracts should be the same.
// This function stores the code for that contract only once, but it can
// be instantiated with custom inputs in the future.
//
// TODO: return gas cost? Add gas limit??? there is no metering here...
func (w *Wasmer) Create(code WasmCode) (CodeID, error) {
	return api.Create(w.cache, code)
}

// GetCode will load the original wasm code for the given code id.
// This will only succeed if that code id was previously returned from
// a call to Create.
//
// This can be used so that the (short) code id (hash) is stored in the iavl tree
// and the larger binary blobs (wasm and pre-compiles) are all managed by the
// rust library
func (w *Wasmer) GetCode(code CodeID) (WasmCode, error) {
	return api.GetCode(w.cache, code)
}

// Instantiate will create a new contract based on the given codeID.
// We can set the initMsg (contract "genesis") here, and it then receives
// an account and address and can be invoked (Execute) many times.
//
// Storage should be set with a PrefixedKVStore that this code can safely access.
//
// Under the hood, we may recompile the wasm, use a cached native compile, or even use a cached instance
// for performance.
func (w *Wasmer) Instantiate(
	code CodeID,
	env types.Env,
	info types.MessageInfo,
	initMsg []byte,
	store KVStore,
	goapi GoAPI,
	querier Querier,
	gasMeter GasMeter,
	gasLimit uint64,
) (*types.InitResponse, uint64, error) {
	envBin, err := json.Marshal(env)
	if err != nil {
		return nil, 0, err
	}
	infoBin, err := json.Marshal(info)
	if err != nil {
		return nil, 0, err
	}
	data, gasUsed, err := api.Instantiate(w.cache, code, envBin, infoBin, initMsg, &gasMeter, store, &goapi, &querier, gasLimit)
	if err != nil {
		return nil, gasUsed, err
	}

	var resp types.InitResult
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, gasUsed, err
	}
	if resp.Err != "" {
		return nil, gasUsed, fmt.Errorf("%s", resp.Err)
	}
	return resp.Ok, gasUsed, nil
}

// Execute calls a given contract. Since the only difference between contracts with the same CodeID is the
// data in their local storage, and their address in the outside world, we need no ContractID here.
// (That is a detail for the external, sdk-facing, side).
//
// The caller is responsible for passing the correct `store` (which must have been initialized exactly once),
// and setting the env with relevent info on this instance (address, balance, etc)
func (w *Wasmer) Execute(
	code CodeID,
	env types.Env,
	info types.MessageInfo,
	executeMsg []byte,
	store KVStore,
	goapi GoAPI,
	querier Querier,
	gasMeter GasMeter,
	gasLimit uint64,
) (*types.HandleResponse, uint64, error) {
	envBin, err := json.Marshal(env)
	if err != nil {
		return nil, 0, err
	}
	infoBin, err := json.Marshal(info)
	if err != nil {
		return nil, 0, err
	}
	data, gasUsed, err := api.Handle(w.cache, code, envBin, infoBin, executeMsg, &gasMeter, store, &goapi, &querier, gasLimit)
	if err != nil {
		return nil, gasUsed, err
	}

	var resp types.HandleResult
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, gasUsed, err
	}
	if resp.Err != "" {
		return nil, gasUsed, fmt.Errorf("%s", resp.Err)
	}
	return resp.Ok, gasUsed, nil
}

// Query allows a client to execute a contract-specific query. If the result is not empty, it should be
// valid json-encoded data to return to the client.
// The meaning of path and data can be determined by the code. Path is the suffix of the abci.QueryRequest.Path
func (w *Wasmer) Query(
	code CodeID,
	env types.Env,
	queryMsg []byte,
	store KVStore,
	goapi GoAPI,
	querier Querier,
	gasMeter GasMeter,
	gasLimit uint64,
) ([]byte, uint64, error) {
	envBin, err := json.Marshal(env)
	if err != nil {
		return nil, 0, err
	}
	data, gasUsed, err := api.Query(w.cache, code, envBin, queryMsg, &gasMeter, store, &goapi, &querier, gasLimit)
	if err != nil {
		return nil, gasUsed, err
	}

	var resp types.QueryResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, gasUsed, err
	}
	if resp.Err != "" {
		return nil, gasUsed, fmt.Errorf("%s", resp.Err)
	}
	return resp.Ok, gasUsed, nil
}

// Migrate will migrate an existing contract to a new code binary.
// This takes storage of the data from the original contract and the CodeID of the new contract that should
// replace it. This allows it to run a migration step if needed, or return an error if unable to migrate
// the given data.
//
// MigrateMsg has some data on how to perform the migration.
func (w *Wasmer) Migrate(
	code CodeID,
	env types.Env,
	info types.MessageInfo,
	migrateMsg []byte,
	store KVStore,
	goapi GoAPI,
	querier Querier,
	gasMeter GasMeter,
	gasLimit uint64,
) (*types.MigrateResponse, uint64, error) {
	envBin, err := json.Marshal(env)
	if err != nil {
		return nil, 0, err
	}
	infoBin, err := json.Marshal(info)
	if err != nil {
		return nil, 0, err
	}
	data, gasUsed, err := api.Migrate(w.cache, code, envBin, infoBin, migrateMsg, &gasMeter, store, &goapi, &querier, gasLimit)
	if err != nil {
		return nil, gasUsed, err
	}

	var resp types.MigrateResult
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, gasUsed, err
	}
	if resp.Err != "" {
		return nil, gasUsed, fmt.Errorf("%s", resp.Err)
	}
	return resp.Ok, gasUsed, nil
}