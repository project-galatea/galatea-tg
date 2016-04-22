package main

import (
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/cli"
)

var (
	BotToken  string
	Version   string
	BuildTime string
)

func main() {
	app := cli.NewApp()

	app.Name = "Galatea Telegram"
	app.Usage = "The main server, load balancer, and telegram interface for the Galatea Project"

	app.Authors = []cli.Author{
		cli.Author{
			Name: "Aidan Lloyd-Tucker",
		},
		cli.Author{
			Name: "Leif",
		},
		cli.Author{
			Name: "Nachi",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "c, connections",
			Usage:  "List of slaves for the master to connect to",
			EnvVar: "CONN_SLAVES",
		},
	}

	app.Version = Version

	num, err := strconv.ParseInt(BuildTime, 10, 64)
	if err == nil {
		app.Compiled = time.Unix(num, 0)
	}

	app.Action = runApp
	app.Run(os.Args)
}

func runApp(c *cli.Context) {
	for _, ip := range c.StringSlice("c") {
		ConnectNewFollower(ip)
	}
	startBot(BotToken)
}
