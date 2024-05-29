package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var nodeList = []string{
	"/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT",
}

// If using vscode, go to Preferences > Settings, and edit Go: Test Timeout to at least 60s

func parseTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}

func getValue(param string) string {
	x := strings.Split(param, "=")
	return x[1]
}

func getIntValue(param string) int64 {
	x := strings.Split(param, "=")

	num, err := strconv.ParseInt(x[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func getArrValue(param string) []string {
	x := strings.ReplaceAll(param, "\"", "")
	x = strings.ReplaceAll(x, "[", "")
	x = strings.ReplaceAll(x, "]", "")

	x2 := strings.Split(x, " ")
	return x2
}

func (s *StoreSuite) TestBasic() {
	ctx := context.Background()

	// Connecting to nodes
	// ================================================================

	log.Info("Connecting to nodes...")

	connectToNodes(ctx, s.node)

	time.Sleep(2 * time.Second) // Required so Identify protocol is executed

	// Open the file
	file, err := os.Open("missing_messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Store
	// ================================================================

	// Read each line and process it

	contentTopc := make(map[string]struct{})

	for scanner.Scan() {
		line := scanner.Text()

		params := strings.Split(line, ";")

		startTime := time.Unix(0, getIntValue(params[1]))
		endTime := time.Unix(0, getIntValue(params[2]))
		contentTopics := getArrValue(getValue(params[3]))
		pubsubTopic := getValue(params[4])
		for _, x := range contentTopics {
			contentTopc[x] = struct{}{}
		}

		cnt, err := queryNode(ctx, s.node, nodeList[0], pubsubTopic, contentTopics, startTime, endTime)
		if err != nil {
			fmt.Println("COULD NOT QUERY STORENODE: ", err)
		} else {
			fmt.Println(cnt, "MESSAGES FOUND FOR - ", startTime, endTime, (contentTopics), pubsubTopic)
			if cnt != 0 {
				fmt.Println("!!!!!!!!!!!!!!!!")
			}
		}

	}

}
