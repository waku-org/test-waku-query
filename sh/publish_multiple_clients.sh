#!/bin/bash

## This script allows to publish multiple messages from multiple processes.

## Params
## $1 - The first parameter passed represents the number of clients.
## $2 - The second parameter has the # of messages each client will publish.

for i in $(seq $1); do
  bash sh/publish_one_client.sh $2 &
  pids[${i}]=$!
done

# wait for all pids
for pid in ${pids[*]}; do
  echo Waiting for ${pid}
  wait $pid
done
