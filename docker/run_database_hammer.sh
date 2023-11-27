#!/bin/sh

## This script is aimed to generate traffic to the database and simulate
## simulate multiple Waku nodes connected to the same database.

## Install postgresql-client for psql command
apk --no-cache add postgresql-client

## The port numbers and host names are defined in the 'docker-compose.yml' file

echo "This is database hammer"

export PGPASSWORD=test123
QUERY_LAST_FIVE_MINUTES=$(echo 'SELECT * FROM messages ORDER BY storedAt DESC LIMIT 100;')

while true
do
  echo "Before making query"
  psql -h postgres -p 5432 -U postgres -U postgres -w -c "${QUERY_LAST_FIVE_MINUTES}"

  sleep 0.01
done
