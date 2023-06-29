package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"go.uber.org/zap"
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

func sendMessages(ctx context.Context, node *node.WakuNode, numMsgToSend int, topic string, contentTopic string) error {
	for i := 0; i < numMsgToSend; i++ {
		payload := make([]byte, 128)
		_, err := rand.Read(payload)
		if err != nil {
			return err
		}

		msg := &pb.WakuMessage{
			Payload:      payload,
			Version:      0,
			ContentTopic: contentTopic,
			Timestamp:    node.Timesource().Now().UnixNano(),
		}

		_, err = node.Relay().Publish(ctx, msg)
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func queryNode(ctx context.Context, node *node.WakuNode, addr string, pubsubTopic string, contentTopic string, startTime time.Time, endTime time.Time) (int, error) {
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
		ContentTopics: []string{contentTopic},
		StartTime:     startTime.UnixNano(),
		EndTime:       endTime.UnixNano(),
	}, store.WithPeer(info.ID), store.WithPaging(false, 100), store.WithRequestId([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
	if err != nil {
		return -1, err
	}

	for {
		cursorIterations += 1
		hasNext, err := result.Next(ctx)
		if err != nil {
			return -1, err
		}

		if !hasNext { // No more messages available
			break
		}

		cnt += len(result.GetMessages())
	}

	log.Info(fmt.Sprintf("%d messages found in %s (Used cursor %d times)\n", cnt, info.ID, cursorIterations))

	return cnt, nil
}
