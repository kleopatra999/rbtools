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

var target_host string

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
	}

	app.Commands = []cli.Command{
		commands.Update,
	}

	app.Run(os.Args)
}
