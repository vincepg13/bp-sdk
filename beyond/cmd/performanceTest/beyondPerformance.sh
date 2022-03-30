#!/bin/bash
nodes=( ["1"]="172.31.16.154:26657"
        ["2"]="172.31.25.77:26657"
        ["3"]="172.31.31.135:26657"
        ["4"]="172.31.26.184:26657")

start=$1
end=$2

for ((a=$start; a <= $end ; a++))
do
randNode=`grep -m1 -ao '[1-4]' /dev/urandom | sed s/0/4/ | head -n1` 
./executeTx.sh car"$a" station"$a" /home/ubuntu/.beyondcli-"$a"/ ${nodes[$randNode]} &
sleep 2
done

