package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pso-dev/qatiq/backend/authentication-service/internal/service"
)

const port = 80

var counts int64

type configuration struct {
	port int64
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
}

type api struct {
	cfg     configuration
	service *service.Service
	wg      sync.WaitGroup
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var app api
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	app.cfg.port = *flags.Int64("port", port, "port to listen on")
	app.cfg.env = *flags.String("env", "development", "environment [development | production]")
	app.cfg.db.dsn = *flags.String("db-dsn", "postgres://postgres:password@localhost/qatiq?sslmode-disable", "data source name")
	app.cfg.db.maxOpenConns = *flags.Int("db-max-open-conns", 25, "database maximum open connections")
	app.cfg.db.maxIdleConns = *flags.Int("db-max-idle-conns", 25, "database maximum idle connections")
	app.cfg.db.maxIdleTime = *flags.String("db-max-idle-time", "15m", "database maximum idle time")
	flags.Func("trusted-origins", "trusted origins (space separated list)", func(val string) error {
		app.cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	app.service = &service.Service{Name: "authentication-service", DB: db}
	return app.serve()
}

func connect() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			counts++
		} else {
			return connection, nil
		}
		if counts > 10 {
			return nil, err
		}
		time.Sleep(2 * time.Second)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// db.SetMaxOpenConns(maxOpenConns)
	// db.SetMaxIdleConns(maxIdleConns)
	// db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
