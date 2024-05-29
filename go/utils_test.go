package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol/legacy_store"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func connectToNodes(ctx context.Context, node *node.WakuNode) {
	wg := sync.WaitGroup{}
	for _, addr := range nodeList {
		wg.Add(1)
		go func(addr string) {
			wg.Done()
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			err := node.DialPeer(ctx, addr)
			if err != nil {
				log.Error("could not connect to peer", zap.String("addr", addr), zap.Error(err))
			}
		}(addr)
	}
	wg.Wait()
}

func queryNode(ctx context.Context, node *node.WakuNode, addr string, pubsubTopic string, contentTopics []string, startTime time.Time, endTime time.Time) (int, error) {
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

	result, err := node.LegacyStore().Query(ctx, legacy_store.Query{
		PubsubTopic:   pubsubTopic,
		ContentTopics: contentTopics,
		StartTime:     proto.Int64(startTime.UnixNano()),
		EndTime:       proto.Int64(endTime.UnixNano()),
	}, legacy_store.WithPeer(info.ID), legacy_store.WithPaging(false, 100))
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

		cnt += len(result.GetMessages())
	}

	log.Info(fmt.Sprintf("%d messages found in %s (Used cursor %d times)\n", cnt, info.ID, cursorIterations))

	return cnt, nil
}
