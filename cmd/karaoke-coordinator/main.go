package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/myl7/karaoke/pkg/coordinator"
)

func main() {
	opts := parseOpts()
	ctx := context.Background()

	co := coordinator.NewCoordinator(coordinator.CoordinatorConfig{
		ServerN: opts["sn"].(int),
		RAddr:   opts["raddr"].(string),
		MURI:    opts["muri"].(string),
	})
	err := co.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := co.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()
	log.Println("bootstrap OK")

	err = co.Run(ctx)
	if err != nil {
		panic(err)
	}
}

func parseOpts() map[string]any {
	sn := flag.Int("sn", 0, "server num")
	raddr := flag.String("raddr", "", "Redis addr")
	muri := flag.String("muri", "", "MongoDB URI")
	flag.Parse()

	abort := func() {
		flag.Usage()
		os.Exit(1)
	}
	if *sn <= 0 {
		abort()
	}
	if *raddr == "" {
		abort()
	}
	if *muri == "" {
		abort()
	}

	opts := make(map[string]any)
	opts["sn"] = *sn
	opts["raddr"] = *raddr
	opts["muri"] = *muri
	return opts
}
