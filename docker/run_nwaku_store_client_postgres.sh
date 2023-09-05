#!/bin/sh

## Install the `dig` command
apk add --update bind-tools

peer_IP=$(dig +short nwaku_postgres)

exec /usr/bin/wakunode\
  --nodekey=7d714a1fada214dead6dc9c7274585eca0ff292451866e7d6d677dc818e8ccd2\
  --storenode=/ip4/${peer_IP}/tcp/30303/p2p/16Uiu2HAmJyLCRhiErTRFcW5GKPrpoMjGbbWdFMx4GCUnnhmxeYhd\
  --log-level=ERROR\
  --rest=true\
  --rest-port=8645\
  --rest-address=0.0.0.0
