package main

import (
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"

	"github.com/status-im/status-go/sdk"
)

type argT struct {
	cli.Helper
	Username string `cli:"username" usage:"Username of the bot account" dft:"the-spectator"`
	Password string `cli:"password" usage:"Password of the bot account" dft:"the-spectator-pwd"`
	Channel  string `cli:"channel" usage:"Channel that bot listens to" dft:"humans-need-not-apply"`
	Interval int    `cli:"interval" usage:"Send message every x second" dft:"5"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {

		rpcClient, err := rpc.Dial("http://localhost:8545")
		checkErr(err)

		client := sdk.New(rpcClient)

		a := client.Readonly()

		r := gin.Default()
		r.LoadHTMLGlob("_assets/html/*")

		monitorChannel := func(channel string) *messagesStore {
			messages := NewMessagesStore(10000, channel)

			ch, err := a.JoinPublicChannel(channel)
			checkErr(err)

			_, _ = ch.Subscribe(func(m *sdk.Msg) {
				if m != nil {
					log.Println("Message from ", m.From, " with body: ", m.Raw)
					if err := messages.Add(*m); err != nil {
						log.Printf("Error while storing message: ERR: %v", err)
					}
				} else {
					log.Println("received a nil message!")
				}
			})

			r.GET("/"+channel, func(c *gin.Context) {
				c.HTML(http.StatusOK, "index.tmpl", gin.H{
					"ChannelName": channel,
					"Messages":    validMessages(messages.Messages()),
				})
			})

			return messages
		}

		for _, channel := range []string{"ethdenver", "teambuilding", "status", "test"} {
			db := monitorChannel(channel)
			defer db.Close()
		}

		r.Run() // listen and serve on 0.0.0.0:8080

		return nil
	})
}

func validMessages(msgs []Msg) []Msg {
	result := []Msg{}

	for _, candidate := range msgs {
		if candidate.Valid() {
			result = append(result, candidate)
		}
	}

	return result
}
