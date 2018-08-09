// pinger.go
// Public chat bot that sends a message every x seconds

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"

	"github.com/ethereum/go-ethereum/rpc"
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

		rpcClient, err := rpc.Dial("http://localhost:8545")

		if err != nil {
			panic("couldn't connect to the local statusd instance")
		}

		client := sdk.New(rpcClient)

		user, err := client.SignupAndLogin(conf.Password)

		if err != nil {
			panic(err)
		}

		pingFunc := func(c *gin.Context, user *sdk.Account, channel string, interval, count int) {
			ch, err := user.JoinPublicChannel(channel)
			if err != nil {
				panic("Couldn't connect to channel: " + channel + "reason: " + err.Error())
			}

			messagesSent := 0
			for range time.Tick(time.Duration(interval) * time.Millisecond) {
				message := fmt.Sprintf("PING no %d : %d", messagesSent, time.Now().Unix())
				err := ch.Publish(message)
				if err != nil {
					fmt.Println("error while publishing -> " + err.Error())
				}
				messagesSent++

				c.Writer.WriteString(fmt.Sprintf("* SENT:  %17d   MESSAGES -> %s*\n", messagesSent, message))
				c.Writer.Flush()

				if messagesSent >= count {
					break
				}
			}

			c.Writer.WriteString("DONE")
			c.Writer.Flush()
		}

		r := gin.Default()
		// ping as the default user
		r.GET("/ping/:channel", func(c *gin.Context) {
			interval := MustParseIntFromQuery(c, "interval", "1000")
			count := MustParseIntFromQuery(c, "count", "1")
			channel := c.Param("channel")
			c.Writer.WriteHeader(200)
			pingFunc(c, user, channel, interval, count)
			c.Writer.CloseNotify()
		})

		// ping as a new user for every request
		r.GET("/ping-as-user/:channel", func(c *gin.Context) {
			user, err := client.SignupAndLogin(conf.Password)
			if err != nil {
				panic(err)
			}
			interval := MustParseIntFromQuery(c, "interval", "1000")
			count := MustParseIntFromQuery(c, "count", "1")
			channel := c.Param("channel")
			c.Writer.WriteHeader(200)
			pingFunc(c, user, channel, interval, count)
			c.Writer.CloseNotify()
		})

		// makes a stress test with N users sending M messages each with an interval
		r.GET("/stress-test/:channel", func(c *gin.Context) {
			interval := MustParseIntFromQuery(c, "interval", "1000")
			count := MustParseIntFromQuery(c, "count", "1")
			usersCount := MustParseIntFromQuery(c, "users", "1")
			channel := c.Param("channel")
			c.Writer.WriteHeader(200)
			wg := &sync.WaitGroup{}
			for i := 0; i < usersCount; i++ {
				wg.Add(1)
				var userNo = i
				time.Sleep(500 * time.Millisecond)
				go func() {
					log := fmt.Sprintf("user %d started", userNo)
					fmt.Println(log)
					c.Writer.WriteString(log)
					c.Writer.Flush()

					user, err := client.SignupAndLogin(conf.Password)
					if err != nil {
						panic(err)
					}

					pingFunc(c, user, channel, interval, count)

					log = fmt.Sprintf("user %d done", userNo)
					fmt.Println(log)
					c.Writer.WriteString(log)
					c.Writer.Flush()

					wg.Done()
				}()
			}

			wg.Wait()

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
