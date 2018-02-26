package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"

	"github.com/mandrigin/status-go-bots/bots"
)

func main() {
	cli.Run(&bots.Config{}, func(ctx *cli.Context) error {
		conf := ctx.Argv().(*bots.Config)

		messages := NewMessagesStore(1000)
		defer messages.Close()

		log.Println("conf: ", conf)

		node := bots.Quickstart(*conf, 100*time.Millisecond, func(ch *bots.StatusChannel) {
			for _, msg := range ch.ReadMessages() {
				if err := messages.Add(msg); err != nil {
					log.Printf("Error while storing message: ERR: %v", err)
				}
			}
		})

		log.Println("Node started, %v", node)

		r := gin.Default()
		r.LoadHTMLGlob("_assets/html/*")
		r.GET("/json", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"messages": messages.Messages(conf.Channel),
			})
		})
		r.GET("/html", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"Messages": messages.Messages(conf.Channel),
			})
		})
		r.Run() // listen and serve on 0.0.0.0:8080

		return nil
	})
}
