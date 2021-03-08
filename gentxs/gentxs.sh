#!/bin/bash
sender=$1
recipient=$2
n_txs=$3
number=`nftchaincli query account $sender | jq .value.account_number | sed 's/"//g'`
sequence=`nftchaincli query account $sender | jq .value.sequence | sed 's/"//g'`
#echo $number
#echo $sequence
nftchaincli "tx" "send" $sender $recipient "1token" "--generate-only" "--from" $sender "-y" > tx_$sender.json
# iterate 0 to tx number - 1
for i in `seq 0 1 $(($n_txs-1))`
	
do
	nftchaincli tx sign tx_$sender.json "--account-number" $number "--sequence" $((sequence+i)) "--offline" "--from" $sender > signed_tx_$sender.json
	nftchaincli tx encode signed_tx_$sender.json >> tmp_$sender.txt
done
sed -e 's/\"//g' tmp_$sender.txt > ./txs/txs_$sender.txt
rm tmp_$sender.txt
rm tx_$sender.json
rm signed_tx_$sender.json
echo done $sender
