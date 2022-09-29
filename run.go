package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/status-im/go-waku/waku/v2/protocol/pb"
	"github.com/status-im/go-waku/waku/v2/protocol/store"
	"github.com/status-im/go-waku/waku/v2/utils"
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host1, err := libp2p.New(DefaultLibP2POptions...)
	if err != nil {
		panic(err)
	}

	s1 := store.NewWakuStore(host1, nil, nil, utils.Logger())
	s1.Start(ctx)
	defer s1.Stop()

	p, err := multiaddr.NewMultiaddr("/dns4/node-02.do-ams3.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmSve7tR5YZugpskMv2dmJAsMUKmfWYEKRXNUxRaTCnsXV")
	if err != nil {
		panic(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(p)
	if err != nil {
		panic(err)
	}

	err = host1.Connect(ctx, *info)
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)

	// Cursor for first query:
	c1 := &pb.Index{
		Digest:       []byte{188, 194, 250, 11, 122, 9, 225, 198, 229, 26, 55, 76, 35, 21, 32, 89, 138, 224, 220, 79, 160, 27, 63, 239, 182, 158, 89, 19, 79, 54, 132, 38},
		ReceiverTime: 1664302928000000000,
		SenderTime:   1664302928245493000,
		PubsubTopic:  "/waku/2/default-waku/proto",
	}

	result, err := s1.Query(ctx, store.Query{
		Topic:         "/waku/2/default-waku/proto",
		ContentTopics: []string{"/waku/1/0x53278eae/rfc26", "/waku/1/0xdfe5d73f/rfc26", "/waku/1/0xb5f68014/rfc26", "/waku/1/0xa0309f58/rfc26", "/waku/1/0x41238e4e/rfc26", "/waku/1/0xfacce554/rfc26", "/waku/1/0x5712b01d/rfc26", "/waku/1/0xca21b94e/rfc26", "/waku/1/0xf702b427/rfc26", "/waku/1/0x657a859c/rfc26"},
		StartTime:     1664302919000000000,
		EndTime:       1664460418000000000,
	}, store.WithPeer(info.ID), store.WithPaging(false, 1000), store.WithCursor(c1))
	if err != nil {
		panic(err)
	}

	fmt.Println("MESSAGES RETRIEVED IN FIRST QUERY:", len(result.Messages))

	fmt.Println("SAME CURSOR?", c1.PubsubTopic == result.Cursor().PubsubTopic && c1.ReceiverTime == result.Cursor().ReceiverTime && c1.SenderTime == result.Cursor().SenderTime && bytes.Equal(c1.Digest, result.Cursor().Digest))

	result2, err := s1.Query(ctx, store.Query{
		Topic:         "/waku/2/default-waku/proto",
		ContentTopics: []string{"/waku/1/0x53278eae/rfc26", "/waku/1/0xdfe5d73f/rfc26", "/waku/1/0xb5f68014/rfc26", "/waku/1/0xa0309f58/rfc26", "/waku/1/0x41238e4e/rfc26", "/waku/1/0xfacce554/rfc26", "/waku/1/0x5712b01d/rfc26", "/waku/1/0xca21b94e/rfc26", "/waku/1/0xf702b427/rfc26", "/waku/1/0x657a859c/rfc26"},
		StartTime:     1664302919000000000,
		EndTime:       1664460418000000000,
	}, store.WithPeer(info.ID), store.WithPaging(false, 1000), store.WithCursor(result.Cursor()))
	if err != nil {
		panic(err)
	}

	fmt.Println("MESSAGES RETRIEVED IN SECOND QUERY:", len(result2.Messages))

}
