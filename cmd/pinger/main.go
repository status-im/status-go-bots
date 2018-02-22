package main

import (
	"fmt"
	"time"

	"github.com/mkideal/cli"

	"github.com/mandrigin/status-go-bots/bots"
)

func main() {
	cli.Run(&bots.Config{}, func(ctx *cli.Context) error {
		conf := ctx.Argv().(*bots.Config)

		node := bots.Quickstart(conf, 10*time.Second, func(ch *bots.StatusChannel) {
			message := fmt.Sprintf("Gopher, gopher: %d", time.Now().Unix())
			ch.SendMessage(message)
		})

		// wait till node has been stopped
		node.Wait()
	}
}
