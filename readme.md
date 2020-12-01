# nftchain

**nftchain** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Get started

```
starport serve
```

`serve` command installs dependencies, initializes and runs the application.

## Configure

Initialization parameters of your app are stored in `config.yml`.

### `accounts`

A list of user accounts created during genesis of your application.

| Key   | Required | Type            | Description                                       |
| ----- | -------- | --------------- | ------------------------------------------------- |
| name  | Y        | String          | Local name of the key pair                        |
| coins | Y        | List of Strings | Initial coins with denominations (e.g. "100coin") |

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

### Contract


