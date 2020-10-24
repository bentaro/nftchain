rm -r $HOME/.nftexchanged/config
rm -r $HOME/.nftexchanged/data
rm -r $HOME/.nftexchangecli/keyring-test-cosmos
nft-exchanged unsafe-reset-all
nft-exchanged init mynode --chain-id nft

nft-exchangecli config keyring-backend test
nft-exchangecli config chain-id nft
nft-exchangecli config output json
nft-exchangecli config indent true
nft-exchangecli config trust-node true

nft-exchangecli keys add user1
nft-exchangecli keys add user2
nft-exchanged add-genesis-account $(nft-exchangecli keys show user1 -a) 1000token,100000000stake
nft-exchanged add-genesis-account $(nft-exchangecli keys show user2 -a) 500token

nft-exchanged gentx --name user1 --keyring-backend test

nft-exchanged collect-gentxs
