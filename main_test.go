package main

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
	// TODO: search criteria
	pubsubTopic := relay.DefaultWakuTopic
	contentTopics := []string{"test1"}
	envelopeHash := "0x" // Use "0x" to find all messages that match the pubsub topic, content topic and start/end time
	startTime := time.Now().Add(-20 * time.Second)
	endTime := time.Now()

	// =========================================================

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Test shouldnt take more than 60s
	defer cancel()

	addNodes(ctx, s.node)
	hash, err := hexutil.Decode(envelopeHash)
	if err != nil {
		panic("invalid envelope hash id")
	}

	wg := sync.WaitGroup{}
	for _, addr := range nodeList {
		wg.Add(1)
		func(addr string) {
			defer wg.Done()
			_, err := queryNode(ctx, s.node, addr, pubsubTopic, contentTopics, startTime, endTime, hash)
			s.NoError(err)
		}(addr)
	}
	wg.Wait()
}
