package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) {
		println("boom! I say!")
	}

	app.Run(os.Args)
}
