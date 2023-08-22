
### Summary

This simple script stress the _Store_ protocol from a different approach.

The setup needs to have two [`nwaku`](https://github.com/waku-org/nwaku)
running, A & B.

Content of "cfg_node_a.txt" file:
```code
ports-shift = 1
pubsub-topic = [ "/waku/2/default-waku/proto" "/waku/2/testing-store" "/waku/2/dev-waku/proto" ]
staticnode = [ "/ip4/0.0.0.0/tcp/60000/p2p/16Uiu2HAmVFXtAfSj4EiR7mL2KvL4EE2wztuQgUSBoj2Jx2KeXFLN" ]
storenode = "/ip4/127.0.0.1/tcp/60000/p2p/16Uiu2HAmVFXtAfSj4EiR7mL2KvL4EE2wztuQgUSBoj2Jx2KeXFLN"
log-level = "DEBUG"
nodekey = "364d111d729a6eb6d2e6113e163f017b5ef03a6f94c9b5b7bb1bb36fa5cb07a9"
rest = true
lightpush = true
discv5-discovery = true
discv5-udp-port = 9000
discv5-enr-auto-update = false
rpc-admin = true
metrics-server = true
```

Content of "cfg_node_b.txt" file:
```code
ports-shift = 0
pubsub-topic = [ "/waku/2/default-waku/proto" "/waku/2/testing-store" ]
staticnode = [ "/ip4/0.0.0.0/tcp/60001/p2p/16Uiu2HAm2eqzqp6xn32fzgGi8K4BuF88W4Xy6yxsmDcW8h1gj6ie" ]
log-level = "DEBUG"
nodekey = "0d714a1fada214dead6dc9c7274585eca0ff292451866e7d6d677dc818e8ccd2"
lightpush = true
store = true
store-message-db-url = "postgres://postgres:test123@localhost:5432/postgres"
#store-message-db-url = "sqlite://sqlite_folder/store.sqlite3"
store-message-retention-policy = "time:6000"
rpc-admin=true
metrics-server = true
```

### Setup

#### Start a Postgres database
`docker compose -f postgres-docker-compose.yml up -d`

#### Run node A
1. Open one terminal and go to the root folder of the [`nwaku`](https://github.com/waku-org/nwaku) repo.
2. Run `./build/wakunode2 --config-file=cfg_node_a.txt`

#### Run node B
1. Open one terminal and go to the root folder of the [`nwaku`](https://github.com/waku-org/nwaku) repo.
2. Run `./build/wakunode2 --config-file=cfg_node_b.txt`

( notice that node B is connected to a database )

#### Send messages to node A

The next example will start 25 processes that each will send 300 messages.

`bash sh/publish_multiple_clients.sh 25 300`
