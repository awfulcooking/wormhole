package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"awful.cooking/wormhole"
)

var staticDir = flag.String("staticDir", "", "directory from which to serve static web files. if omitted, an embedded copy of the controller app will be served.")

func main() {
	flag.Parse()

	config := wormhole.DefaultServerConfig()
	config.NameGenerator = HumanNameGenerator
	config.WebsocketReadLimit = 100 * 1024 * 1024
	if *staticDir != "" {
		config.StaticFS = os.DirFS(*staticDir)
	}

	if len(flag.Args()) == 0 {
		log.Fatal("please provide an address to listen on as the first argument")
	}

	if err := Run(flag.Arg(0), config); err != nil {
		log.Fatal(err)
	}
}

// Run starts a http.Server for the passed in address
// with all requests handled by a wormhole.WebsocketHandler
func Run(addr string, config wormhole.ServerConfig) error {
	routes := http.NewServeMux()
	routes.Handle("/", wormhole.NewRouter(config))

	server := &http.Server{
		Handler:           routes,
		ReadHeaderTimeout: time.Second * 10,
	}

	if l, err := net.Listen("tcp", addr); err != nil {
		return err
	} else {
		log.Printf("listening on http://%v", l.Addr())
		go server.Serve(l)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("Shutting down: %v", <-interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return server.Shutdown(ctx)
}
