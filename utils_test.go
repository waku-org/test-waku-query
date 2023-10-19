package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
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
			panic(err)
		}

		msg := &pb.WakuMessage{
			Payload:      payload,
			Version:      0,
			ContentTopic: contentTopic,
			Timestamp:    node.Timesource().Now().UnixNano(),
		}

		_, err = node.Relay().Publish(ctx, msg)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func sendMessagesConcurrent(ctx context.Context, node *node.WakuNode, numMsgToSend int, topic string, contentTopic string) error {
	wg := sync.WaitGroup{}
	for i := 0; i < numMsgToSend; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			payload := make([]byte, 128)
			_, err := rand.Read(payload)
			if err != nil {
				panic(err)
			}

			msg := &pb.WakuMessage{
				Payload:      payload,
				Version:      0,
				ContentTopic: contentTopic,
				Timestamp:    node.Timesource().Now().UnixNano(),
			}

			_, err = node.Relay().Publish(ctx, msg)
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
	return nil
}

func queryNode(ctx context.Context, node *node.WakuNode, addr string, pubsubTopic string, startTime time.Time, endTime time.Time) (int, error) {
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
		Topic:         "/waku/2/default-waku/proto",
		StartTime:     time.Now().Add(time.Duration(-30) * time.Minute).UnixNano(),
		EndTime:       time.Now().UnixNano(),
		ContentTopics: []string{"/waku/1/0xee3a5ba0/rfc26"},
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

		for _, msg := range result.GetMessages() {
			env := protocol.NewEnvelope(msg, time.Now().UnixNano(), relay.DefaultWakuTopic)

			envHash := hexutil.Encode(env.Hash())
			fmt.Println(envHash, env.Message().ContentTopic, env.Message().Timestamp)
		}

		cnt += len(result.GetMessages())
	}

	log.Info(fmt.Sprintf("%d messages found in %s (Used cursor %d times)\n", cnt, info.ID, cursorIterations))

	return cnt, nil
}
