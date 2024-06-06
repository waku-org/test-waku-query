# Query tool

Use this to query a storenode (uses StoreV2)
```
make

./build/query \
  --cluster-id=16 \
  --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
  --pubsub-topic=/waku/2/rs/16/32 \
  --content-topic=/waku/1/0x242ed557/rfc26 \
  --content-topic=/waku/1/0xd811cd50/rfc26 \
  --content-topic=/waku/1/0x89bed93d/rfc26 \
  --content-topic=/waku/1/0xc95d2429/rfc26 \
  --content-topic=/waku/1/0xa0a6b41b/rfc26 \
  --start-time=1717507412000000000 \
  --end-time=1717593812000000000
```

For the previous execution, you should see among the logs the following:
```
2024-06-06T23:23:56.840Z        INFO    query   go/main.go:228  TOTAL MESSAGES RETRIEVED        {"node": "16Uiu2HAmNtg3tBjhaGyjQosPh4AitwLKJtcLwoVC5ycTuf4BCh7V", "num": 30}
```

### Docker
```
# Build
docker build -t querytool:latest .

# Execute
docker run querytool:latest \
  --cluster-id=16 \
  --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
  --pubsub-topic=/waku/2/rs/16/32 \
  --content-topic=/waku/1/0x242ed557/rfc26 \
  --content-topic=/waku/1/0xd811cd50/rfc26 \
  --content-topic=/waku/1/0x89bed93d/rfc26 \
  --content-topic=/waku/1/0xc95d2429/rfc26 \
  --content-topic=/waku/1/0xa0a6b41b/rfc26 \
  --start-time=1717507412000000000 \
  --end-time=1717593812000000000
```
