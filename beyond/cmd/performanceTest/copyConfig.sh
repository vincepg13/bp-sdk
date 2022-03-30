#!/bin/bash
for ((a=1; a <= 100 ; a++))
do
    echo "preparing folder:"  beyondcli-"$a" 
    rm -rf .beyondcli-"$a"
    mkdir -p .beyondcli-"$a"/keys
    cp -r keys.db .beyondcli-"$a"/keys/
done
