#!/bin/sh

## Install the `dig` command
apk add --update bind-tools

peer_IP=$(dig +short nwaku_sqlite)

exec /usr/bin/wakunode\
  --nodekey=5d714a1fada214dead6dc9c7274585eca0ff292451866e7d6d677dc818e8ccd2\
  --storenode=/ip4/${peer_IP}/tcp/30304/p2p/16Uiu2HAkxj3WzLiqBximSaHc8wV9Co87GyRGRYLVGsHZrzi3TL5W\
  --log-level=ERROR\
  --rest=true\
  --rest-port=8646\
  --rest-address=0.0.0.0
