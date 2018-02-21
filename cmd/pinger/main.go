package main

import (
	"fmt"
	"time"

	"github.com/mandrigin/status-go-bots/bots"
)

func main() {
	conf := bots.Config{Password: "my-cool-password", Channel: "humans-need-not-apply", Nickname: "Cloudy Test Sender "}
	node := bots.Quickstart(conf, 10*time.Second, func(ch *bots.StatusChannel) {
		message := fmt.Sprintf("Gopher, gopher: %d", time.Now().Unix())
		ch.SendMessage(message)
	})

	// wait till node has been stopped
	node.Wait()
}
