package main

import (
	"fmt"
	"context"
	"sync"
	"time"
	"strconv"

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

func parseTime(s string) (time.Time) {
	i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        panic(err)
    }
    return time.Unix(i, 0)
}

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

	// err := sendMessages(  to send the msgs sequentially
	err := sendMessagesConcurrent(ctx, s.node, numMsgToSend, pubsubTopic, contentTopic)
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

func (s *StoreSuite) TestCompareDatabasesPerformance() {
	// The next settings might be to be adapted depending on the databases content.
	// We seek to pick times windows so that the number of returned rows is ~1000.
	expectedNumMsgs := 966
	pubsubTopic := "/waku/2/default-waku/proto"
	startTime := parseTime("1695992040")
	endTime :=   parseTime("1695992056")
	contentTopic := "/waku/2/default-content/proto"

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Test shouldnt take more than 60s
	defer cancel()

	// Connecting to nodes
	// ================================================================

	log.Info("Connecting to nodes...")

	connectToNodes(ctx, s.node)

	time.Sleep(2 * time.Second) // Required so Identify protocol is executed

	// Store
	// ================================================================

    timeSpentMap := make(map[string]time.Duration)
	numUsers := int64(10)

	peers := []string{
		// Postgres peer
		"/ip4/127.0.0.1/tcp/30303/p2p/16Uiu2HAmJyLCRhiErTRFcW5GKPrpoMjGbbWdFMx4GCUnnhmxeYhd",
		// SQLite peer
		"/ip4/127.0.0.1/tcp/30304/p2p/16Uiu2HAkxj3WzLiqBximSaHc8wV9Co87GyRGRYLVGsHZrzi3TL5W",
	}

	wg := sync.WaitGroup{}
	for _, addr := range peers {
		for userIndex := 0; userIndex < int(numUsers); userIndex++ {
			wg.Add(1)
			go func(addr string) {
				defer wg.Done()
				fmt.Println("Querying node", addr)
				start := time.Now()
				cnt, err := queryNode(ctx, s.node, addr, pubsubTopic, contentTopic, startTime, endTime)
				timeSpent := time.Since(start)
				fmt.Printf("\n%s took %v\n\n", addr, timeSpent)
				s.NoError(err)
				s.Equal(expectedNumMsgs, cnt)
				timeSpentMap[addr] += timeSpent
			}(addr)
		}
	}

	wg.Wait()

	timeSpentNanos := timeSpentMap[peers[0]].Nanoseconds() / numUsers
	fmt.Println("\n\nAverage time spent: ", peers[0], time.Duration(timeSpentNanos))

	timeSpentNanos = timeSpentMap[peers[1]].Nanoseconds() / numUsers
	fmt.Println("\n\nAverage time spent:", peers[1], time.Duration(timeSpentNanos))
	fmt.Println("")
}

