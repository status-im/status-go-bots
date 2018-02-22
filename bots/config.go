package bots

import (
	"os"

	"github.com/mkideal/cli"
)

type Config struct {
	cli.Helper
	Channel  string `cli:"channel" usage:"Channel that bot listens to" dft:"humans-need-not-apply"`
	Nickname string `cli:"nickname" usage:"Nickname of the bot" dft:"generic bot"`
}

func (c Config) Password() string {
	return os.Getenv("STATUS_BOT_PASSWORD")
}
