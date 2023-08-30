#!/bin/sh

## Install the `dig` command
apk --no-cache add bind-tools

apk --no-cache add curl

## Install coreutils so that the `date` command behaves as the GNU one
apk --no-cache add coreutils

## The port numbers and host names are defined in the 'docker-compose.yml' file
nwaku_postgres_IP=$(dig +short nwaku_postgres)
nwaku_sqlite_IP=$(dig +short nwaku_sqlite)
target_postgres="http://${nwaku_postgres_IP}:8545"
target_sqlite="http://${nwaku_sqlite_IP}:8546"

echo "This is publisher: ${target_postgres}; ${target_sqlite}"

## Wait a few seconds until the `nwaku` nodes started their rpc services
sleep 20

while true
do
  ## Send a 'get_waku_v2_store_v1_messages' req to the ""postgres"" node
  curl -d '{"jsonrpc":"2.0","id":"id","method":"get_waku_v2_store_v1_messages"}' --header "Content-Type: application/json" ${target_postgres}

  ## Send a 'get_waku_v2_store_v1_messages' req to the ""sqlite"" node
  curl -d '{"jsonrpc":"2.0","id":"id","method":"get_waku_v2_store_v1_messages"}' --header "Content-Type: application/json" ${target_sqlite}

  sleep 5
done
