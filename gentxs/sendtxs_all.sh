#!/bin/bash
host=$1
senders=($(cat senders|xargs))
if [ -z "$host" ]; then
	exit
fi
ids=()

for sender in ${senders[@]}
do 
	number=`nftchaincli query account $sender | jq .value.account_number | sed 's/"//g'`
        ids=("${ids[@]}" $number)
done

# iterate 0 to length of senders - 1
for i in `seq 0 1 $((${#ids[@]}-1))`

do
	./sendtxs.sh $host ${ids[$i]} ${senders[$i]} &	
done

echo

