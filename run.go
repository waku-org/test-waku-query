package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/status-im/status-go/eth-node/types"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"github.com/waku-org/go-waku/waku/v2/timesource"
	"github.com/waku-org/go-waku/waku/v2/utils"
)

// Default options used in the libp2p node
var DefaultLibP2POptions = []libp2p.Option{
	libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
	),
	libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	),
	libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
}

func ContentTopic(t []byte) string {
	enc := hexutil.Encode(t)
	return "/waku/1/" + enc + "/rfc26"
}

// ToTopic converts a string to a whisper topic.
func ToTopic(s string) []byte {
	return crypto.Keccak256([]byte(s))[:types.TopicLength]
}

func buildComparerFN(messages []string) (func([]byte) bool, error) {
	var msgIDArray [][]byte
	for _, m := range messages {
		b, err := hexutil.Decode(m)
		if err != nil {
			return nil, err
		}
		msgIDArray = append(msgIDArray, b)
	}

	return func(hash []byte) bool {
		for _, m := range msgIDArray {
			if bytes.Equal(m, hash) {
				return true
			}
		}
		return false
	}, nil
}

var nodeList = []string{
	"/dns4/node-01.ac-cn-hongkong-c.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAkvEZgh3KLwhLwXg95e5ojM8XykJ4Kxi2T7hk22rnA7pJC",
	"/dns4/node-01.do-ams3.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAm6HZZr7aToTvEBPpiys4UxajCTU97zj5v7RNR2gbniy1D",
	"/dns4/node-01.gc-us-central1-a.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAkwBp8T6G77kQXSNMnxgaMky1JeyML5yqoTHRM8dbeCBNb",
	"/dns4/node-02.ac-cn-hongkong-c.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmFy8BrJhCEmCYrUfBdSNkrPw6VHExtv4rRp1DSBnCPgx8",
	"/dns4/node-02.do-ams3.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmSve7tR5YZugpskMv2dmJAsMUKmfWYEKRXNUxRaTCnsXV",
	"/dns4/node-02.gc-us-central1-a.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmDQugwDHM3YeUp86iGjrUvbdw3JPRgikC7YoGBsT2ymMg",
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host1, err := libp2p.New(DefaultLibP2POptions...)
	if err != nil {
		panic(err)
	}

	s1 := store.NewWakuStore(host1, nil, nil, timesource.NewDefaultClock(), utils.Logger())
	s1.Start(ctx)
	defer s1.Stop()

	for i, n := range nodeList {
		queryNode(ctx, n, host1, s1, i)
	}
}

func queryNode(ctx context.Context, node string, host1 host.Host, s1 *store.WakuStore, i int) {
	p, err := multiaddr.NewMultiaddr(node)
	if err != nil {
		panic(err)
	}

	messagesToLookFor := []string{
		"0x305b9f88bc8f670b57d895b06296a50eef8b69f1576bce7313d61c3fd4adf677",
		"0x5336a19ad110eb5efaa39c180881529da490e3cc12b6f634d00a4105ed57da21",
		"0x7554ffbeb8ec0373b165013708214ebda8103c9ff3d1bd676e0c1cd9b1a9571e",
	}

	startDate := 1671058980
	endDate := 1671058980 + 120

	messageExists, err := buildComparerFN(messagesToLookFor)
	if err != nil {
		panic(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(p)
	if err != nil {
		panic(err)
	}

	err = host1.Connect(ctx, *info)
	if err != nil {
		fmt.Printf("Could not connect to %s: %s", info.ID, err.Error())
		return
	}

	cnt := 0
	cursorIterations := 0

	result, err := s1.Query(ctx, store.Query{
		Topic:     "/waku/2/default-waku/proto",
		StartTime: int64(time.Duration(startDate) * time.Second),
		EndTime:   int64((time.Duration(endDate) * time.Second)),
	}, store.WithPeer(info.ID), store.WithPaging(false, 100), store.WithRequestId([]byte{1, 2, 3, 4, 5, 6, 7, 8, byte(i)}))
	if err != nil {
		fmt.Printf("Could not query %s: %s", info.ID, err.Error())
		return
	}

	for {
		cnt += len(result.Messages)
		cursorIterations += 1

		for _, r := range result.Messages {
			h, _, _ := r.Hash()
			if messageExists(h) {
				fmt.Println("Message found in", node)
				return
			}
		}

		if result.IsComplete() {
			break
		}

		result, err = s1.Next(ctx, result)
		if err != nil {
			fmt.Printf("Could not retrieve more results from %s: %s", info.ID, err.Error())
		}
	}

	fmt.Println("Message NOT found in", node)

}
