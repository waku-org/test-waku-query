package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/waku-org/go-waku/logging"
	"github.com/waku-org/go-waku/waku/cliutils"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/legacy_store"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"github.com/waku-org/go-waku/waku/v2/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

type Options struct {
	NodeKey       *ecdsa.PrivateKey
	ClusterID     uint
	PubSubTopic   string
	ContentTopics cli.StringSlice
	StartTime     int64
	EndTime       int64
	PageSize      uint64
	StoreNode     *multiaddr.Multiaddr
	UseLegacy     bool
	QueryTimeout  time.Duration
	LogLevel      string
	LogEncoding   string
	LogOutput     string
}

var options Options

var flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{Name: "config-file", Usage: "loads configuration from a TOML file (cmd-line parameters take precedence)"},
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:  "nodekey",
		Usage: "P2P node private key as hex.",
		Value: &cliutils.PrivateKeyValue{
			Value: &options.NodeKey,
		},
	}),
	altsrc.NewUintFlag(&cli.UintFlag{
		Name:        "cluster-id",
		Value:       0,
		Usage:       "Cluster id that the node is running in. Node in a different cluster id is disconnected.",
		Destination: &options.ClusterID,
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "pubsub-topic",
		Usage:       "Query pubsub topic.",
		Destination: &options.PubSubTopic,
		Required:    true,
	}),
	altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:        "content-topic",
		Usage:       "Query content topic. Argument may be repeated.",
		Destination: &options.ContentTopics,
	}),
	altsrc.NewInt64Flag(&cli.Int64Flag{
		Name:        "start-time",
		Usage:       "Query start time in nanoseconds",
		Destination: &options.StartTime,
	}),
	altsrc.NewInt64Flag(&cli.Int64Flag{
		Name:        "end-time",
		Usage:       "Query end time in nanoseconds",
		Destination: &options.EndTime,
	}),
	altsrc.NewUint64Flag(&cli.Uint64Flag{
		Name:        "pagesize",
		Value:       20,
		Usage:       "Pagesize",
		Destination: &options.PageSize,
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:  "storenode",
		Usage: "Multiaddr of a peer that supports store protocol",
		Value: &cliutils.MultiaddrValue{
			Value: &options.StoreNode,
		},
		Required: true,
	}),
	altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:        "use-legacy",
		Usage:       "Use legacy store",
		Destination: &options.UseLegacy,
	}),
	altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:        "timeout",
		Usage:       "timeout for each individual store query request",
		Destination: &options.QueryTimeout,
		Value:       1 * time.Minute,
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Value: &cliutils.ChoiceValue{
			Choices: []string{"DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", "FATAL"},
			Value:   &options.LogLevel,
		},
		Usage: "Define the logging level (allowed values: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL)",
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:  "log-encoding",
		Usage: "Define the encoding used for the logs (allowed values: console, nocolor, json)",
		Value: &cliutils.ChoiceValue{
			Choices: []string{"console", "nocolor", "json"},
			Value:   &options.LogEncoding,
		},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "log-output",
		Value:       "stdout",
		Usage:       "specifies where logging output should be written  (stdout, file, file:./filename.log)",
		Destination: &options.LogOutput,
	}),
}

func main() {
	// Defaults
	options.LogLevel = "INFO"
	options.LogEncoding = "console"

	app := &cli.App{
		Name:    "query",
		Version: "0.0.1",
		Before:  altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc("config-file")),
		Flags:   flags,
		Action: func(c *cli.Context) error {
			err := Execute(c.Context, options)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func Execute(ctx context.Context, opts Options) error {
	utils.InitLogger(options.LogEncoding, options.LogOutput, "query")

	var prvKey *ecdsa.PrivateKey
	var err error

	if options.NodeKey != nil {
		prvKey = options.NodeKey
	} else {
		if prvKey, err = crypto.GenerateKey(); err != nil {
			return fmt.Errorf("error generating key: %w", err)
		}
	}

	p2pPrvKey := utils.EcdsaPrivKeyToSecp256k1PrivKey(prvKey)
	id, err := peer.IDFromPublicKey(p2pPrvKey.GetPublic())
	if err != nil {
		return err
	}
	logger := utils.Logger().With(logging.HostID("node", id))

	lvl, err := zapcore.ParseLevel(options.LogLevel)
	if err != nil {
		return err
	}

	libp2pOpts := append(node.DefaultLibP2POptions, libp2p.NATPortMap()) // Attempt to open ports using uPNP for NATed hosts.)

	wakuNode, err := node.New(
		node.WithLogger(logger),
		node.WithLogLevel(lvl),
		node.WithPrivateKey(prvKey),
		node.WithClusterID(uint16(options.ClusterID)),
		node.WithNTP(),
		node.WithLibP2POptions(libp2pOpts...),
	)
	if err != nil {
		return fmt.Errorf("could not instantiate waku: %w", err)
	}

	if err = wakuNode.Start(ctx); err != nil {
		return err
	}
	defer wakuNode.Stop()

	cnt := 0

	if !options.UseLegacy {
		criteria := store.FilterCriteria{
			ContentFilter: protocol.NewContentFilter(options.PubSubTopic, options.ContentTopics.Value()...),
			TimeStart:     proto.Int64(options.StartTime),
			TimeEnd:       proto.Int64(options.EndTime),
		}

		ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
		result, err := wakuNode.Store().Query(ctx, criteria,
			store.WithPeerAddr(*options.StoreNode),
			store.WithPaging(false, options.PageSize),
			store.IncludeData(false),
		)
		cancel()
		if err != nil {
			return err
		}

		for !result.IsComplete() {
			cnt += len(result.Messages())

			ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
			err := result.Next(ctx)
			cancel()
			if err != nil {
				return err
			}
		}

	} else {
		query := legacy_store.Query{
			PubsubTopic:   options.PubSubTopic,
			ContentTopics: options.ContentTopics.Value(),
			StartTime:     proto.Int64(options.StartTime),
			EndTime:       proto.Int64(options.EndTime),
		}

		ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
		result, err := wakuNode.LegacyStore().Query(ctx, query,
			legacy_store.WithPeerAddr(*options.StoreNode),
			legacy_store.WithPaging(false, 20),
		)
		cancel()
		if err != nil {
			return err
		}

		for {
			ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
			hasNext, err := result.Next(ctx)
			cancel()
			if err != nil {
				return err
			}

			if !hasNext { // No more messages available
				break
			}

			cnt += len(result.GetMessages())
		}
	}

	logger.Info("TOTAL MESSAGES RETRIEVED", zap.Int("num", cnt))

	return nil
}
