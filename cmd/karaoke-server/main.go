package main

import (
	"context"
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/myl7/karaoke/pkg/server"
)

func main() {
	opts := parseOpts()
	ctx := context.Background()

	s := server.NewServer(server.ServerConfig{
		Layer: opts["layer"].(int),
		Addr:  opts["addr"].(string),
		LAddr: opts["laddr"].(string),
		RAddr: opts["raddr"].(string),
		MURI:  opts["muri"].(string),
	})
	go func() {
		err := s.Listen(ctx)
		if err != nil {
			panic(err)
		}
	}()
	err := s.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := s.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()
	log.Println("bootstrap OK")

	err = s.Run(ctx)
	if err != nil {
		panic(err)
	}
}

func parseOpts() map[string]any {
	layer := flag.Int("layer", 0, "layer in Karaoke")
	addr := flag.String("addr", "", "endpoint accessiable by other servers")
	laddr := flag.String("laddr", "", "gRPC listen addr")
	raddr := flag.String("raddr", "", "Redis addr")
	muri := flag.String("muri", "", "MongoDB URI")
	flag.Parse()

	abort := func() {
		flag.Usage()
		os.Exit(1)
	}
	if *layer < 0 {
		abort()
	}
	if *laddr == "" {
		abort()
	}
	if *addr == "" && !regexp.MustCompile(`:\d+`).Match([]byte(*laddr)) {
		abort()
	}
	if *raddr == "" {
		abort()
	}
	if *muri == "" {
		abort()
	}

	opts := make(map[string]any)
	opts["layer"] = *layer
	opts["addr"] = *addr
	opts["laddr"] = *laddr
	opts["raddr"] = *raddr
	opts["muri"] = *muri
	return opts
}
