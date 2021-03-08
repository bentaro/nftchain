#!/bin/bash
senders=($(cat senders|xargs))
recipient=$1
txs_num=$2
if [ -z "$recipient" ]; then
	exit
fi
ids=()
for sender in ${senders[@]}
do
	number=`nftchaincli query account $sender | jq .value.account_number | sed 's/"//g'`
	ids=($number "${ids[@]}")
done
for sender in ${senders[@]}
do
	./gentxs.sh $sender $recipient $txs_num &
done


