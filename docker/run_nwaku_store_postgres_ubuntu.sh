#!/bin/sh

apt-get update
## Install the `dig` command
apt-get install dnsutils -y
apt-get install wget -y

bootstrap_IP=$(dig +short bootstrap)

apt-get install libpq5 -y
chmod +x /usr/bin/wakunode

RETRIES=${RETRIES:=10}

while [ -z "${BOOTSTRAP_ENR}" ] && [ ${RETRIES} -ge 0 ]; do
  BOOTSTRAP_ENR=$(wget -O - --post-data='{"jsonrpc":"2.0","method":"get_waku_v2_debug_v1_info","params":[],"id":1}' --header='Content-Type:application/json' http://${bootstrap_IP}:8544/ 2> /dev/null | sed 's/.*"enrUri":"\([^"]*\)".*/\1/');
  echo "Bootstrap node not ready in ${bootstrap_IP}, retrying (retries left: ${RETRIES})"
  sleep 1
  RETRIES=$(( $RETRIES - 1 ))
done

if [ -z "${BOOTSTRAP_ENR}" ]; then
   echo "Could not get BOOTSTRAP_ENR and none provided. Failing"
   exit 1
fi

IP=$(hostname -I)

echo "I am postgres ubuntu. Listening on: ${IP}"

./usr/bin/wakunode\
  --relay=true\
  --topic=/waku/2/default-waku/proto\
  --topic=/waku/2/dev-waku/proto\
  --rpc-admin=true\
  --keep-alive=true\
  --log-level=DEBUG\
  --rpc-port=8545\
  --rpc-address=0.0.0.0\
  --metrics-server=True\
  --metrics-server-port=8003\
  --metrics-server-address=0.0.0.0\
  --max-connections=4\
  --dns-discovery=true\
  --discv5-discovery=true\
  --discv5-enr-auto-update=True\
  --discv5-bootstrap-node=${BOOTSTRAP_ENR}\
  --nat=extip:${IP}\
  --store-message-db-url="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/postgres"\
  --store=true\
  --store-message-retention-policy=capacity:12000000
