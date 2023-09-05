
## Summary
This repo contains tools for analysing performance around the Store protocol

### Golang tool

The Golang project is aimed to setup `go-waku` clients that publish messages
and later make requests trough the _Store_ protocol to retrieve the stored
messages

To run the tests, go to the `go` folder and run the `make` command.
Notice that the `go` tool expects a running `nwaku` node(s) to be running with
_Store_ protocol mounted, and a running postgres database.

### Bash tool

Simple script that allows to publish messages from different clients.

`BASH`(n clients) --json-rpc--> `nwaku_A` <--relay--> `nwaku_B` <---> `database`

Notice that the bash script expects two `nwaku` nodes that communicate through
the _Relay_ protocol and the `nwaku_B` has the _Store_ protocol mounted and
is connected to the `postgres_DB`.

### Docker

Contains a docker compose file with:
- Two `nwaku` nodes configured with _Postgres_ and _SQLite_.
- Grafana container to compare performance of both nodes.
- Container with simple shell script that sends publish requests through rpc.
- Two `nwaku` nodes configured as _Store_-clients and listening to REST requests.

Inside the _docker_ folder, we have a _jmeter_ test plan which is aimed for
performing concurrent _Store_ REST requests to the _Store_-clients.
