package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pso-dev/qatiq/backend/logger-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const port = 80

type configuration struct {
	port int64
	env  string
	db   struct {
		url string
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
	app.cfg.db.url = *flags.String("client", "mongodb://mongo:27017", "mongo client url")
	flags.Func("trusted-origins", "trusted origins (space separated list)", func(val string) error {
		app.cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	client, err := connectToMongo(app.cfg.db.url)
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err) // add recover panic middleware
		}
	}()
	app.service = &service.Service{DB: client}
	return app.serve()
}

func connectToMongo(url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",    // TODO: remove magic
		Password: "password", // TODO: remove magic
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting:", err)
		return nil, err
	}
	return c, nil
}
