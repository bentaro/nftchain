# nftchain

**nftchain** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Get started

```
nftchaind start
```
## CLI commands

### General

Show all accounts in local device

```nftchaincli keys list```

Query account balance

```nftchaincli query account <account address>```

Send token

```nftchaincli tx send <sender address> <recipient address> <amount||denom(ex. 100token)>```

Query transaction detail

```nftchaincli query tx <transaction hash>```
 
### NFT

Mint NFT

```nftchaincli tx nft mint <nft denom> <NFT address> <recipient address> --from <transactor address>```

Transfer NFT

```nftchaincli tx nft transfer <sender address> <recipient address> <nft denom> <NFT address> --from <transactor address>```

Query NFTs by owner

```nftchaincli query nft <owner address>```

Show all denoms in Network

```nftchaincli query nft denoms```
 
Query token info

```nftchaincli query nft token <denom> <NFT address>```

### Contract

Deploy contract to the network 

```nftchaincli tx wasm store contract.wasm --from <transactor address>```

Instantiate contract

```nftchaincli tx wasm instantiate <contract index> <constructor members (ex. "{\"denom\": \"token\"}") > --from <transactor address> --label <unique label string>```


Query contract

```nftchaincli query wasm contract-state smart <contract address> <query json (ex. "{\"config\": {}}") >```

Execute Contract

```nftchaincli tx wasm execute <contract address> <query json (ex. "{\"create_poll\": {\"quorum_percentage\": 0,\"description\": \"first poll\",\"start_height\": 1,\"end_height\": 6000}}")> --from <transactor address>```
