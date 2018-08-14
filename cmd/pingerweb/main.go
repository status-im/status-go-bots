// pinger.go
// Public chat bot that sends a message every x seconds

package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
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

	quotes := []string{
		"Perfection is achieved, not when there is nothing more to add, but when there is nothing left to take away.\nAntoine de Saint-Exupery",
		"I think there is a profound and enduring beauty in simplicity; in clarity, in efficiency. True simplicity is derived from so much more than just the absence of clutter and ornamentation. It's about bringing order to complexity.\n Jonathan Ive",
		"Simplification is one of the most difficult things to do. <https://www.azquotes.com/quote/780210>\nJonathan Ive",
		"In many ways, it’s the things that are not there that we are most proud of. For us, it is all about refining and refining until it seems like there’s nothing between the user and the content they are interacting with.\nJonathan Ive",
		"“Confusion and clutter are the failure of design, not the attributes of information.\nEdward R Tufte",
		"'Simple' is a tricky word, it can mean a lot of things. To us, it just means clear. That doesn't always mean total reduction, or minimalism - sometimes, to make things clearer, you have to add a step. <https://www.azquotes.com/quote/1498137>\nJason Fried",
		"In a world where everybody screams, silence is noticeable. White space provides the silence.\n - Vignelli",
		"“Design isn’t crafting a beautiful textured button with breathtaking animation. It’s figuring out if there’s a way to get rid of the button altogether.\n — Edward Tufte",
		"Simple is hard. Easy is harder. Invisible is hardest.\n — Jean-Louis Gassée",
		"The ability to simplify means to eliminate the unnecessary so that the necessary may speak.  —Hans Hofmann <http://www.hanshofmann.org/>",
		"Fools ignore complexity. Pragmatists suffer it. Some can avoid it. Geniuses remove it.  Alan J. Perlis",
		"Simplicity is not the goal. It is the by-product of a good idea and modest expectations. –Paul Rand",
		"Good design is obvious. Great design is transparent.” –Joe Sparano <https://twitter.com/intent/tweet?text=%22%22Good+design+is+obvious.+Great+design+is+transparent.%22+%E2%80%93Joe+Sparano%22https%3A%2F%2Fwww.invisionapp.com%2Fblog%2Fdesign-and-creativity-quotes%2F+via+%40InVisionApp>",
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator

	randomQuote := func() string {
		idx := r.Intn(len(quotes))
		source := quotes[idx]
		reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
		if err != nil {
			log.Fatal(err)
		}
		processedString := reg.ReplaceAllString(source, " ")
		return processedString
	}

	fmt.Println(randomQuote)

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
				message := fmt.Sprintf("PING no %d: %s", messagesSent, randomQuote())
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
				time.Sleep(1000 * time.Millisecond)
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
