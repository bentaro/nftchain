#!/bin/bash
host=$1
id=$2
sender=$3
txs=($(cat ./txs/txs_$sender.txt|xargs))

for i in `seq 0 1 $((${#txs[@]}-1))`
do
	json=$(cat << EOS 
	{"jsonrpc":"2.0","id":${id},"method":"broadcast_tx_sync","params":{"tx":"${txs[i]}"}}
EOS
)
curl -s -XPOST $host -H "Content-Type: application/json" -d $json > /dev/null
done

echo done $sender

