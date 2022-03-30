for pid in $(ps -ef | awk '/executeTx/ {print $2}'); do kill -9 $pid; done
