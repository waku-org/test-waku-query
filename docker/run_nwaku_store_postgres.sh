#!/bin/sh

## Install the `dig` command
apk add --update bind-tools

peer_IP=$(dig +short nwaku_sqlite)

exec /usr/bin/wakunode\
  --nodekey=1d714a1fada214dead6dc9c7274585eca0ff292451866e7d6d677dc818e8ccd2\
  --staticnode=/ip4/${peer_IP}/tcp/30304/p2p/16Uiu2HAkxj3WzLiqBximSaHc8wV9Co87GyRGRYLVGsHZrzi3TL5W\
  --relay=true\
  --topic=/waku/2/default-waku/proto\
  --topic=/waku/2/dev-waku/proto\
  --rpc-admin=true\
  --keep-alive=true\
  --log-level=ERROR\
  --rpc-port=8545\
  --rpc-address=0.0.0.0\
  --tcp-port=30303\
  --metrics-server=True\
  --metrics-server-port=8003\
  --metrics-server-address=0.0.0.0\
  --store-message-db-url="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/postgres"\
  --store=true\
  --store-message-retention-policy=capacity:4000000
