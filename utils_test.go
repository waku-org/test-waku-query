package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/waku-org/go-waku/logging"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/peers"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"go.uber.org/zap"
)

func addNodes(ctx context.Context, node *node.WakuNode) {
	for _, addr := range nodeList {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Error("invalid multiaddress", zap.Error(err), zap.String("addr", addr))
			continue
		}

		_, err = node.AddPeer(ma, peers.Static, store.StoreID_v20beta4)
		if err != nil {
			log.Error("could not add peer", zap.Error(err), zap.Stringer("addr", ma))
			continue
		}
	}
}

func queryNode(ctx context.Context, node *node.WakuNode, addr string, pubsubTopic string, contentTopics []string, startTime time.Time, endTime time.Time, envelopeHash []byte) (int, error) {
	p, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return -1, err
	}

	info, err := peer.AddrInfoFromP2pAddr(p)
	if err != nil {
		return -1, err
	}

	cnt := 0
	cursorIterations := 0

	result, err := node.Store().Query(ctx, store.Query{
		Topic:         pubsubTopic,
		ContentTopics: contentTopics,
		StartTime:     startTime.UnixNano(),
		EndTime:       endTime.UnixNano(),
	}, store.WithPeer(info.ID), store.WithPaging(false, 100), store.WithRequestId([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
	if err != nil {
		return -1, err
	}

	for {
		hasNext, err := result.Next(ctx)
		if err != nil {
			return -1, err
		}

		if !hasNext { // No more messages available
			break
		}
		cursorIterations += 1

		// uncomment to find message by ID
		for _, m := range result.GetMessages() {
			if len(envelopeHash) != 0 && bytes.Equal(m.Hash(pubsubTopic), envelopeHash) {
				log.Info("MESSAGE FOUND!", logging.HexString("envelopeHash", envelopeHash), logging.HostID("peerID", info.ID))
				return 0, nil
			}
		}

		cnt += len(result.GetMessages())
	}

	log.Info(fmt.Sprintf("%d messages found in %s (Used cursor %d times)\n", cnt, info.ID, cursorIterations))

	return cnt, nil
}
