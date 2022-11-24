package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "jasn port scanner",
		Version: "v0.1",
		Usage:   "scan the port in given ip range, example -i 192.0.2.0/24 -p 22,80-139",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "iplist",
				Aliases: []string{"i"},
				Usage:   "ip list",
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "port list",
				Value:   "",
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "timeout",
				Value:   3,
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"c"},
				Usage:   "concurrency",
				Value:   1000,
			},
		},
		Action: doScan,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doScan(ctx *cli.Context) error {
	log.Println("do scan now")
	return nil
}
