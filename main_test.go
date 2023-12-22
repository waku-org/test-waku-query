package main

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var nodeList = []string{
	"/dns4/node-01.do-ams3.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAm6HZZr7aToTvEBPpiys4UxajCTU97zj5v7RNR2gbniy1D",
}

// If using vscode, go to Preferences > Settings, and edit Go: Test Timeout to at least 60s

func (s *StoreSuite) TestBasic() {
	// TODO: search criteria
	pubsubTopic := "/waku/2/default-waku/proto"
	contentTopics := []string{"/waku/1/0xee3a5ba0/rfc26"}
	envelopeHash := "0x" // Use "0x" to find all messages that match the pubsub topic, content topic and start/end time
	startTime := time.Unix(0, 1700576989000000000)
	endTime := time.Unix(0, 1703255389000000000)

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
