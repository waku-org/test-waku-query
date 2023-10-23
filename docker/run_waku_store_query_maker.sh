#!/bin/sh

IP=$(ip a | grep "inet " | grep -Fv 127.0.0.1 | sed 's/.*inet \([^/]*\).*/\1/')

echo "I am a waku store query generator"

## Getting the address of the Postgres node
RETRIES=10
while [ -z "${POSTGRES_ADDR}" ] && [ ${RETRIES} -ge 0 ]; do
  POSTGRES_ADDR=$(wget -O - --post-data='{"jsonrpc":"2.0","method":"get_waku_v2_debug_v1_info","params":[],"id":1}' --header='Content-Type:application/json' http://nwaku-postgres:8545/ 2> /dev/null | sed 's/.*"listenAddresses":\["\(.*\)"].*/\1/');
  wget -O - --post-data='{"jsonrpc":"2.0","method":"get_waku_v2_debug_v1_info","params":[],"id":1}' --header='Content-Type:application/json' http://nwaku-postgres:8545/
  echo "Store Postgres node node not ready, retrying (retries left: ${RETRIES})"
  sleep 1
  RETRIES=$(( $RETRIES - 1 ))
done

if [ -z "${POSTGRES_ADDR}" ]; then
   echo "Could not get POSTGRES_ADDR and none provided. Failing"
   exit 1
fi

## Getting the address of the SQLite node
RETRIES=10
while [ -z "${SQLITE_ADDR}" ] && [ ${RETRIES} -ge 0 ]; do
  SQLITE_ADDR=$(wget -O - --post-data='{"jsonrpc":"2.0","method":"get_waku_v2_debug_v1_info","params":[],"id":1}' --header='Content-Type:application/json' http://nwaku-sqlite:8546/ 2> /dev/null | sed 's/.*"listenAddresses":\["\(.*\)"].*/\1/');
  echo "Store SQLite node node not ready, retrying (retries left: ${RETRIES})"
  sleep 1
  RETRIES=$(( $RETRIES - 1 ))
done

if [ -z "${SQLITE_ADDR}" ]; then
   echo "Could not get SQLITE_ADDR and none provided. Failing"
   exit 1
fi

## Getting the bootstrap node ENR
RETRIES=10
while [ -z "${BOOTSTRAP_ENR}" ] && [ ${RETRIES} -ge 0 ]; do
  BOOTSTRAP_ENR=$(wget -O - --post-data='{"jsonrpc":"2.0","method":"get_waku_v2_debug_v1_info","params":[],"id":1}' --header='Content-Type:application/json' http://bootstrap:8544/ 2> /dev/null | sed 's/.*"enrUri":"\([^"]*\)".*/\1/');
  echo "Bootstrap node not ready, retrying (retries left: ${RETRIES})"
  sleep 1
  RETRIES=$(( $RETRIES - 1 ))
done

if [ -z "${BOOTSTRAP_ENR}" ]; then
   echo "Could not get BOOTSTRAP_ENR and none provided. Failing"
   exit 1
fi

## Further parameter details:
#     --num-minutes-query -> to indicate the time window in the query. storedAt field.

echo "Using bootstrap node: ${BOOTSTRAP_ENR}"
exec /main\
    --pubsub-topic="/waku/2/default-waku/proto"\
    --content-topic="my-ctopic"\
    --queries-per-second=${STORE_QUERIES_PER_SECOND}\
    --bootstrap-node=${BOOTSTRAP_ENR}\
    --peer-store-postgres-addr="${POSTGRES_ADDR}"\
    --peer-store-sqlite-addr="${SQLITE_ADDR}"\
    --num-minutes-query=60\
    --num-concurrent-users=1
