package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/utils"
)

var log = utils.Logger().Named("TEST")

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

type StoreSuite struct {
	suite.Suite
	nodes []*node.WakuNode
}

func (s *StoreSuite) SetupSuite() {
	for i := 0; i < 5; i++ {
		wakuNode, err := node.New(
			node.WithNTP(),
			node.WithWakuRelayAndMinPeers(1),
		)

		s.NoError(err)

		err = wakuNode.Start(context.Background())
		s.NoError(err)

		s.nodes = append(s.nodes, wakuNode)
	}
}

func (s *StoreSuite) TearDownSuite() {
	for _, x := range s.nodes {
		x.Stop()
	}

}
