// pinger.go
// Public chat bot that sends a message every x seconds

package main

import (
	"fmt"
	"time"

	"github.com/mkideal/cli"
	"github.com/status-im/status-go/sdk"
)

type argT struct {
	cli.Helper
	Username string `cli:"username" usage:"Username of the bot account" dft:"pinger"`
	Password string `cli:"password" usage:"Password of the bot account" dft:"pinger"`
	Channel  string `cli:"channel" usage:"Channel that bot listens to" dft:"humans-need-not-apply"`
	Interval int    `cli:"interval" usage:"Send message every x milliseconds" dft:"5000"`
	Limit    int    `cli:"limit" usage:"Send X messages and quit" dft:"0"`
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		conf := ctx.Argv().(*argT)

		messagesSent := 0

		conn, err := sdk.Connect(conf.Username, conf.Password)
		if err != nil {
			panic("Couldn't connect to status")
		}

		ch, err := conn.Join(conf.Channel)
		if err != nil {
			panic("Couldn't connect to status")
		}

		for range time.Tick(time.Duration(conf.Interval) * time.Millisecond) {
			messagesSent++
			message := fmt.Sprintf("PING no %d : %d", messagesSent, time.Now().Unix())
			ch.Publish(message)
			fmt.Println("***************************************")
			fmt.Printf("* SENT:  %17d   MESSAGES *\n", messagesSent)
			fmt.Println("***************************************")
			if conf.Limit > 0 && messagesSent >= conf.Limit {
				fmt.Println("done. quitting. sending messages first tho. 5 sec...")
				time.Sleep(5 * time.Second)
				return nil
			}
		}

		return nil
	})
}
