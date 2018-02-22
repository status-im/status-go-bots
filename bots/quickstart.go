package bots

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/node"
	sn "github.com/status-im/status-go/geth/node"
)

func Quickstart(conf Config, repeat time.Duration, f func(ch *StatusChannel)) *node.Node {
	config, err := NodeConfig()
	if err != nil {
		log.Fatalf("Making config failed: %v", err)
	}

	nodeManager := sn.NewNodeManager()
	log.Println("Starting node...")
	err = nodeManager.StartNode(config)
	if err != nil {
		log.Fatalf("Node start failed: %v", err)
	}

	node, err := nodeManager.Node()
	if err != nil {
		log.Fatalf("Getting node failed: %v", err)
	}

	SignupOrLogin(nodeManager, conf.Password).Join(conf.Channel, conf.Nickname).RepeatEvery(repeat, f)

	return node
}
