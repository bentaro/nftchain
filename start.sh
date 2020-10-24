rm -r $HOME/.nftchaind/config
rm -r $HOME/.nftchaind/data
rm -r $HOME/.nftchaincli/keyring-test-cosmos
nftchaind unsafe-reset-all
nftchaind init mynode --chain-id nft

nftchaincli config keyring-backend test
nftchaincli config chain-id nft
nftchaincli config output json
nftchaincli config indent true
nftchaincli config trust-node true

nftchaincli keys add user1
nftchaincli keys add user2
nftchaind add-genesis-account $(nftchaincli keys show user1 -a) 1000token,100000000stake
nftchaind add-genesis-account $(nftchaincli keys show user2 -a) 500token

nftchaind gentx --name user1 --keyring-backend test

nftchaind collect-gentxs

