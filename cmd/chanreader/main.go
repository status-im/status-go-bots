package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandrigin/status-go-bots/bots"
)

func main() {
	messages := NewMessagesStore(1000)
	defer messages.Close()

	conf := bots.Config{Password: "my-cool-password", Channel: "humans-need-not-apply", Nickname: "Cloudy Test Baboon"}
	node := bots.Quickstart(conf, 100*time.Millisecond, func(ch *bots.StatusChannel) {
		for _, msg := range ch.ReadMessages() {
			if err := messages.Add(msg); err != nil {
				log.Printf("Error while storing message: ERR: %v", err)
			}
		}
	})

	log.Println("Node started, %v", node)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"messages": messages.Messages("humans"),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
