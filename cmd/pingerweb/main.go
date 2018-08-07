// pinger.go
// Public chat bot that sends a message every x seconds

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"

	sdk "github.com/status-im/status-go-sdk"
)

type argT struct {
	cli.Helper
	Username string `cli:"username" usage:"Username of the bot account" dft:"pinger"`
	Password string `cli:"password" usage:"Password of the bot account" dft:"pinger"`
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		conf := ctx.Argv().(*argT)

		conn, err := sdk.Connect(conf.Username, conf.Password)
		if err != nil {
			panic("Couldn't connect to status")
		}

		r := gin.Default()
		r.GET("/ping/:channel", func(c *gin.Context) {
			interval := MustParseIntFromQuery(c, "interval", "1000")
			count := MustParseIntFromQuery(c, "count", "1")

			ch, err := conn.Join(c.Param("channel"))
			if err != nil {
				panic("Couldn't connect to channel:" + c.Param("channel"))
			}

			c.Writer.WriteHeader(200)

			messagesSent := 0
			for range time.Tick(time.Duration(interval) * time.Millisecond) {
				message := fmt.Sprintf("PING no %d : %d", messagesSent, time.Now().Unix())
				ch.Publish(message)
				messagesSent++

				c.Writer.WriteString(fmt.Sprintf("* SENT:  %17d   MESSAGES -> %s*\n", messagesSent, message))
				c.Writer.Flush()

				if messagesSent >= count {
					break
				}
			}

			c.Writer.WriteString("DONE")
			c.Writer.Flush()
			c.Writer.CloseNotify()
		})

		r.Run() // listen and serve on 0.0.0.0:8080

		return nil
	})
}

func MustParseIntFromQuery(c *gin.Context, q string, defaultValue string) int {
	valueStr := c.DefaultQuery(q, defaultValue)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(value)
}
