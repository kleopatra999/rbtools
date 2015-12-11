package actions

import (
	"github.com/codegangsta/cli"
	"os"
)

func ValidateRequiredArguments(ctx *cli.Context) {
	var (
		host     = ctx.GlobalString("host")
		username = ctx.GlobalString("username")
		password = ctx.GlobalString("password")
	)

	if host == "" {
		println("--host argument is required!")
		os.Exit(1)
	}

	if username == "" {
		println("--username argument is required!")
		os.Exit(1)
	}

	if password == "" {
		println("--password argument is required!")
		os.Exit(1)
	}
}
