package main

import (
	"os"

	"github.com/andefined/twtarget/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "twtarget"
	app.Version = "0.1.0"
	app.Usage = ""
	app.Description = "twtarget (or Twitter Target) is a CLI tool that collects data from twitter API for a given User."

	app.Commands = []cli.Command{
		{
			Name:      "init",
			Usage:     "Initialize a new Target (User). The command will create a target folder and subfolders with the configuration files and collected data.",
			Action:    commands.Init,
			ArgsUsage: "user",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "conf, c",
					Usage: "Configuration `File` (ex: \"default.yaml\")",
				},
			},
		},
		{
			Name:   "fetch",
			Usage:  "Fetch data for a given Target (User). ",
			Action: commands.Fetch,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tweets, t",
					Usage: "Retrieve Tweets",
				},
				cli.BoolFlag{
					Name:  "likes, l",
					Usage: "Retrieve Likes",
				},
				cli.BoolFlag{
					Name:  "friends, r",
					Usage: "Retrieve Friends",
				},
				cli.BoolFlag{
					Name:  "followers, f",
					Usage: "Retrieve Followers",
				},
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Retrieve Everything",
				},
				cli.BoolFlag{
					Name:  "user, u",
					Usage: "Print the current target",
				},
			},
		},
	}

	app.Run(os.Args)
}
