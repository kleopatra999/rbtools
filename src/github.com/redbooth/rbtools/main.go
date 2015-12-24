package main

import (
	"github.com/codegangsta/cli"
	"github.com/redbooth/rbtools/commands"
	"os"
)

const (
	APP_VER = "0.0.1"
	NAME    = "rbtools"
)

var (
	target_host string
	username    string
	password    string
)

// TODO: Allow config file for most flags
func main() {
	app := cli.NewApp()
	app.Name = NAME
	app.Author = "Redbooth Inc"
	app.Email = "private-cloud@redbooth.com"
	app.Version = APP_VER
	app.Usage = "e.g. rbtools backup create"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "",
			Usage:       "Host to target",
			Destination: &target_host},

		cli.StringFlag{
			Name:        "username",
			Value:       "",
			Usage:       "Basic auth username",
			Destination: &username},

		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "Basic auth password",
			Destination: &username},
	}

	app.Action = func(ctx *cli.Context) {
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

	app.Commands = []cli.Command{
		commands.Backup,
	}

	app.Run(os.Args)
}
