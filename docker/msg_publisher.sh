#!/bin/sh

## Install the `dig` command
apk --no-cache add bind-tools

apk --no-cache add curl

## Install coreutils so that the `date` command behaves as the GNU one
apk --no-cache add coreutils

apk --no-cache add openssl

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
  ## Target of ~10kB
  payload_size=$(( $RANDOM % 2000 + 9000 ))
  payload=$(openssl rand -hex ${payload_size} | base64 | tr -d '\n')

  ## Make the ""postgres"" node to publish a message
  curl -d "{\"jsonrpc\":\"2.0\",\"id\":"$(date +%s%N)",\"method\":\"post_waku_v2_relay_v1_message\", \"params\":[\"/waku/2/default-waku/proto\", {\"timestamp\":"$(date +%s%N)", \"payload\":\"${payload}\"}]}" --header "Content-Type: application/json" ${target_postgres}

  ## Make the ""sqlite"" node to publish a message
  curl -d "{\"jsonrpc\":\"2.0\",\"id\":"$(date +%s%N)",\"method\":\"post_waku_v2_relay_v1_message\", \"params\":[\"/waku/2/default-waku/proto\", {\"timestamp\":"$(date +%s%N)", \"payload\":\"${payload}\"}]}" --header "Content-Type: application/json" ${target_sqlite}

  sleep 0.01
done
