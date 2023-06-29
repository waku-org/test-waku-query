package main

import (
	"context"
	"sync"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
)

var nodeList = []string{
	"/dns4/node-01.ac-cn-hongkong-c.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAkvEZgh3KLwhLwXg95e5ojM8XykJ4Kxi2T7hk22rnA7pJC",
	"/dns4/node-01.do-ams3.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAm6HZZr7aToTvEBPpiys4UxajCTU97zj5v7RNR2gbniy1D",
	"/dns4/node-01.gc-us-central1-a.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAkwBp8T6G77kQXSNMnxgaMky1JeyML5yqoTHRM8dbeCBNb",
	"/dns4/node-02.ac-cn-hongkong-c.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmFy8BrJhCEmCYrUfBdSNkrPw6VHExtv4rRp1DSBnCPgx8",
	"/dns4/node-02.do-ams3.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmSve7tR5YZugpskMv2dmJAsMUKmfWYEKRXNUxRaTCnsXV",
	"/dns4/node-02.gc-us-central1-a.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmDQugwDHM3YeUp86iGjrUvbdw3JPRgikC7YoGBsT2ymMg",
}

// If using vscode, go to Preferences > Settings, and edit Go: Test Timeout to at least 60s

func (s *StoreSuite) TestBasic() {
	numMsgToSend := 100
	pubsubTopic := relay.DefaultWakuTopic
	contentTopic := "test1"

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Test shouldnt take more than 60s
	defer cancel()

	// Connecting to nodes
	// ================================================================

	log.Info("Connecting to nodes...")

	connectToNodes(ctx, s.node)

	time.Sleep(2 * time.Second) // Required so Identify protocol is executed

	s.NotZero(len(s.node.Relay().PubSub().ListPeers(relay.DefaultWakuTopic)), "no peers available")

	// Sending messages
	// ================================================================
	startTime := s.node.Timesource().Now()

	err := sendMessages(ctx, s.node, numMsgToSend, pubsubTopic, contentTopic)
	require.NoError(s.T(), err)

	endTime := s.node.Timesource().Now()

	// Store
	// ================================================================

	time.Sleep(5 * time.Second) // Adding a delay to guarantee that messages are inserted (needed with sqlite)

	wg := sync.WaitGroup{}
	for _, addr := range nodeList {
		wg.Add(1)
		func(addr string) {
			defer wg.Done()
			cnt, err := queryNode(ctx, s.node, addr, pubsubTopic, contentTopic, startTime, endTime)
			s.NoError(err)
			s.Equal(numMsgToSend, cnt)
		}(addr)
	}
	wg.Wait()
}
