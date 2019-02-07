package main

import (
	"log"
	"net/http"
	"net/rpc"

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

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		conf := ctx.Argv().(*argT)

		messages := NewMessagesStore(1000)
		defer messages.Close()

		rpcClient, err := rpc.Dial("http://localhost:8545")
		checkErr(err)

		client := sdk.New(rpcClient)

		a, err := client.SignupAndLogin(conf.Password)
		checkErr(err)

		ch, err := a.JoinPublicChannel("igorm-test")
		checkErr(err)

		_, _ = ch.Subscribe(func(m *sdk.Msg) {
			if m != nil {
				log.Println("Message from ", m.From, " with body: ", m.Raw)
				if err := messages.Add(*msg); err != nil {
					log.Printf("Error while storing message: ERR: %v", err)
				}
			} else {
				log.Println("received a nil message!")
			}
		})

		r := gin.Default()
		r.LoadHTMLGlob("_assets/html/*")
		r.GET("/json", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"messages": messages.Messages(),
			})
		})
		r.GET("/html", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"ChannelName": conf.Channel,
				"Messages":    messages.Messages(),
			})
		})
		r.Run() // listen and serve on 0.0.0.0:8080

		return nil
	})
}
