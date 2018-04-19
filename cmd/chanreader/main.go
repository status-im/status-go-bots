package main

import (
	"log"
	"net/http"

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

		conn, err := sdk.Connect(conf.Username, conf.Password)
		if err != nil {
			panic("Couldn't connect to status")
		}

		ch, err := conn.Join(conf.Channel)
		if err != nil {
			panic("Couldn't connect to status")
		}

		ch.Subscribe(func(msg *sdk.Msg) {
			log.Println("received a message", msg)
			if msg != nil {
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
