#!/bin/bash
display_usage() {
	echo "This script must be run with three args, which contain from, to account names and path to config folder"
    echo "Store password in filename: pass"
	echo -e "\nUsage:\n./executeTx.sh car station /home/ubuntu/.beyondcli-1/ \n"
	}

# if less than two arguments supplied, display usage
	if [  $# -le 1 ]
	then
		display_usage
		exit 1
	fi 
   
trap "exit" INT

#set this manually
chain="awesome-beyond-chain"

#read args 
from=$1
to=$2
home=$3
node=$4

fromAddr=`beyondcli keys list --home=$home | grep -w $from | awk '{print $3}'`
toAddr=`beyondcli keys list --home=$home | grep -w $to | awk '{print $3}'`
startSequence=`beyondcli account $fromAddr --trust-node=true --node=$node --chain-id=$chain  | jq -r '.value."BaseAccount".sequence'`
echo $from
echo $fromAddr
echo $startSequence
((startSequence--))
while true;
do
      ((startSequence++))
      echo "initOrder with sequence: " "$startSequence"
        beyondcli initOrder --from=$from --amount=1 --to=$toAddr --async --sequence="$startSequence" --chain-id=$chain --node=$node --home=$home < pass || break
      ((startSequence++))
      echo "finalizeOrder with sequence. " "$startSequence"
       beyondcli finalizeOrder --from=$from --charge=1 --to=$toAddr --async --sequence="$startSequence" --chain-id=$chain --node=$node --home=$home < pass || break
done

