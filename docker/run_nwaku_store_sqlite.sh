#!/bin/sh

## Install the `dig` command
apk add --update bind-tools

peer_IP=$(dig +short nwaku_postgres)

exec /usr/bin/wakunode\
  --nodekey=2d714a1fada214dead6dc9c7274585eca0ff292451866e7d6d677dc818e8ccd2\
  --staticnode=/ip4/${peer_IP}/tcp/30303/p2p/16Uiu2HAmJyLCRhiErTRFcW5GKPrpoMjGbbWdFMx4GCUnnhmxeYhd\
  --relay=true\
  --topic=/waku/2/default-waku/proto\
  --topic=/waku/2/dev-waku/proto\
  --rpc-admin=true\
  --keep-alive=true\
  --log-level=ERROR\
  --rpc-port=8546\
  --rpc-address=0.0.0.0\
  --tcp-port=30304\
  --metrics-server=True\
  --metrics-server-port=8004\
  --metrics-server-address=0.0.0.0\
  --store-message-db-url="sqlite:///data/store.sqlite3"\
  --store=true\
  --store-message-retention-policy=time:864000
