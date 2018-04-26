package main

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/heartbeatsjp/go-ratticcli/commands"
	"gopkg.in/urfave/cli.v1"
)

// Version description
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "rattic"
	app.Usage = "CLI for RatticWeb"
	app.Version = Version

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
			EnvVar: "RATTIC_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "user",
			Value:  myUsername,
			Usage:  "Username",
			EnvVar: "RATTIC_USER",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "RatticWeb TOKEN",
			EnvVar: "RATTIC_TOKEN",
		},
		cli.StringFlag{
			Name:   "cache-path",
			Value:  cachePath,
			Usage:  "Cache File Path",
			EnvVar: "RATTIC_CACHE_PATH",
		},
		cli.IntFlag{
			Name:   "cache-ttl",
			Value:  86400,
			Usage:  "cache ttl(sec)",
			EnvVar: "RATTIC_CACHE_TTL",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "list Cred",
			Action: commands.ListAction,
			Flags:  commands.ListFlags,
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
