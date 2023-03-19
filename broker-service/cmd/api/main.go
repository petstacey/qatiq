package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

const port = 80

type configuration struct {
	port int64
	env  string
	cors struct {
		trustedOrigins []string
	}
}

type broker struct {
	cfg configuration
	wg  sync.WaitGroup
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var app broker
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	app.cfg.port = *flags.Int64("port", port, "port to listen on")
	app.cfg.env = *flags.String("env", "development", "environment [development | production]")
	flags.Func("trusted-origins", "trusted origins (space separated list)", func(val string) error {
		app.cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	app.wg = sync.WaitGroup{}
	return app.serve()
}
