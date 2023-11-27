#!/bin/sh

IP=$(ip a | grep "inet " | grep -Fv 127.0.0.1 | sed 's/.*inet \([^/]*\).*/\1/')

echo "I am a bootstrap node. Listening to: ${IP}"

exec /usr/bin/wakunode\
      --relay=true\
      --peer-exchange=true\
      --peer-persistence=false\
      --max-connections=300\
      --dns-discovery=true\
      --dns-discovery-url=enrtree://AL65EKLJAUXKKPG43HVTML5EFFWEZ7L4LOKTLZCLJASG4DSESQZEC@prod.status.nodes.status.im\
      --discv5-discovery=true\
      --discv5-udp-port=9000\
      --discv5-enr-auto-update=False\
      --rpc-admin=true\
      --keep-alive=true\
      --max-connections=300\
      --dns-discovery=true\
      --discv5-discovery=true\
      --discv5-enr-auto-update=True\
      --log-level=DEBUG\
      --rpc-port=8544\
      --rpc-address=0.0.0.0\
      --metrics-server=True\
      --metrics-server-address=0.0.0.0\
      --nodekey=30348dd51465150e04a5d9d932c72864c8967f806cce60b5d26afeca1e77eb68\
      --nat=extip:${IP}
