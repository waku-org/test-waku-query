#!/bin/sh

## This script is aimed to generate traffic to the database and simulate
## simulate multiple Waku nodes connected to the same database.

## Install the `dig` command
apk --no-cache add bind-tools

## Install postgresql-client for psql command
apk --no-cache add postgresql-client

## The port numbers and host names are defined in the 'docker-compose.yml' file
nwaku-postgres-IP=$(dig +short postgres)
target-postgres="http://${nwaku-postgres-IP}:5432"

echo "This is publisher: ${target-postgres}; ${target_sqlite}"

query-last-five-minutes="SELECT storedAt, contentTopic, payload, pubsubTopic, version, timestamp, id FROM messages WHERE storedAt >= EXTRACT(EPOCH FROM (NOW() - INTERVAL '5 minutes')) * 1000 ORDER BY storedAt DESC;"
while true
do
  echo "Before making query"
  psql -h ${nwaku-postgres-IP} -U postgres -W test123 -U postgres -c "${query-last-five-minutes}"

  sleep 10
done
