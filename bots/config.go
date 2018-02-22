package bots

import "github.com/mkideal/cli"

type Config struct {
	cli.Helper
	Password string `cli:"password" usage:"Password for the account" dft:"my-cool-password"`
	Channel  string `cli:"channel" usage:"Channel that bot listens to" dft:"humans-need-not-apply"`
	Nickname string `cli:"nickname" usage:"Nickname of the bot" dft:"generic bot"`
}
