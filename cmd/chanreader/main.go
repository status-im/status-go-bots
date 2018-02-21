package main

import (
	"log"
	"time"

	"github.com/mandrigin/status-go-bots/bots"
)

func main() {
	messages := NewMessagesStore()
	defer messages.Close()

	conf := bots.Config{Password: "my-cool-password", Channel: "humans-need-not-apply", Nickname: "Cloudy Test Baboon"}
	node := bots.Quickstart(conf, 100*time.Millisecond, func(ch *bots.StatusChannel) {
		for _, msg := range ch.ReadMessages() {
			if err := messages.Add(msg); err != nil {
				log.Printf("Error while storing message: ERR: %v", err)
			}
		}
	})

	go func() {
		for {
			log.Printf("MESSAGES Messages: %d", len(messages.Messages("humans")))
			time.Sleep(1 * time.Second)
		}
	}()

	// wait till node has been stopped
	node.Wait()
}
