package main

import (
	"context"
	"sync"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
)

var nodeList = []string{
	"/dns4/metal-01.he-eu-hel1.vacdev.misc.statusim.net/tcp/60002/p2p/16Uiu2HAmVFXtAfSj4EiR7mL2KvL4EE2wztuQgUSBoj2Jx2KeXFLN",
}

// If using vscode, go to Preferences > Settings, and edit Go: Test Timeout to at least 60s

func (s *StoreSuite) TestBasic() {
	numMsgToSend := 2500
	pubsubTopic := relay.DefaultWakuTopic
	contentTopic := "test22"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second) // Test shouldnt take more than 60s
	defer cancel()

	// Connecting to nodes
	// ================================================================

	log.Info("Connecting to nodes...")

	for _, n := range s.nodes {
		connectToNodes(ctx, n)
	}

	time.Sleep(2 * time.Second) // Required so Identify protocol is executed

	for _, n := range s.nodes {
		s.NotZero(len(n.Relay().PubSub().ListPeers(relay.DefaultWakuTopic)), "no peers available")
	}

	// Sending messages
	// ================================================================
	startTime := s.nodes[0].Timesource().Now().Add(-2 * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(len(s.nodes))
	for _, currN := range s.nodes {
		go func(n *node.WakuNode) {
			defer wg.Done()
			err := sendMessagesConcurrent(ctx, n, numMsgToSend, pubsubTopic, contentTopic)
			require.NoError(s.T(), err)
		}(currN)
	}
	wg.Wait()

	endTime := s.nodes[0].Timesource().Now().Add(2 * time.Second)

	// Store
	// ================================================================

	time.Sleep(5 * time.Second) // Adding a delay to guarantee that messages are inserted (needed with sqlite)

	for _, addr := range nodeList {
		wg.Add(1)
		func(addr string) {
			defer wg.Done()
			cnt, err := queryNode(ctx, s.nodes[0], addr, pubsubTopic, contentTopic, startTime, endTime)
			s.NoError(err)
			s.Equal(numMsgToSend, cnt)
		}(addr)
	}
	wg.Wait()
}
