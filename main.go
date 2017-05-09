package main

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/heartbeatsjp/go-ratticcli/commands"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"

	var myUsername string
	me, err := user.Current()
	if err == nil {
		myUsername = me.Username
	}

	home := os.Getenv("HOME")
	cachePath := filepath.Join(home, ".rattic-cache.db")

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Value:  "https://localhost",
			Usage:  "RatticWeb URL",
			EnvVar: "ENDPOINT",
		},
		cli.StringFlag{
			Name:   "user",
			Value:  myUsername,
			Usage:  "Username",
			EnvVar: "USER",
		},
		cli.StringFlag{
			Name:   "cache-path",
			Value:  cachePath,
			Usage:  "Cache File Path",
			EnvVar: "CACHE_PATH",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "search",
			Usage:  "search Cred",
			Action: commands.SearchAction,
			Flags:  commands.SearchFlags,
		},
		{
			Name:   "show",
			Usage:  "show Cred",
			Action: commands.ShowAction,
			Flags:  commands.ShowFlags,
		},
		{
			Name:   "reload",
			Usage:  "reload token and local cache",
			Action: commands.ReloadAction,
			Flags:  commands.ReloadFlags,
		},
	}
	app.Run(os.Args)
}
