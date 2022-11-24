package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/malfunkt/iprange"
	"github.com/urfave/cli/v2"
)

type scanParams struct {
	ipRange     string
	ipList      []net.IP
	portRange   string
	portList    []int
	timeout     int
	concurrency int
}

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
	sp, err := parseArgs(ctx)
	if err != nil {
		return err
	}

	sp.portList, err = getPorts(sp.portRange)
	if err != nil {
		return err
	}

	sp.ipList, err = getIpList(sp.ipRange)
	if err != nil {
		return err
	}

	log.Println("Begin Scan")
	log.Println(strings.Repeat("*", 100))

	beginScan(sp)

	log.Println(strings.Repeat("*", 100))
	log.Println("End Scan")

	return nil
}

func parseArgs(ctx *cli.Context) (scanParams, error) {
	sp := scanParams{
		portRange:   "22,23,53,80-139",
		timeout:     3,
		concurrency: 1000,
	}

	if ctx.IsSet("iplist") {
		sp.ipRange = ctx.String("iplist")
	} else {
		return sp, errors.New("iplist is required")
	}

	if ctx.IsSet("port") {
		sp.portRange = ctx.String("port")
	}

	if ctx.IsSet("timeout") {
		sp.timeout = ctx.Int("timeout")
	}

	if ctx.IsSet("concurrency") {
		sp.concurrency = ctx.Int("concurrency")
	}

	return sp, nil
}

func getPorts(selection string) ([]int, error) {
	ports := []int{}
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			} else {
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}

func getIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}

	list := addressList.Expand()
	return list, err
}
