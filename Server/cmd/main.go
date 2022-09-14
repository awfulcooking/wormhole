package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"awful.cooking/wormhole"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide an address to listen on as the first argument")
	}

	if err := Run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// Run starts a http.Server for the passed in address
// with all requests handled by a wormhole.WebsocketHandler
func Run(addr string) error {
	config := wormhole.DefaultServerConfig()
	config.NameGenerator = HumanNameGenerator

	routes := http.NewServeMux()
	routes.Handle("/", wormhole.NewRouter(config))

	server := &http.Server{
		Handler:           routes,
		ReadHeaderTimeout: time.Second * 10,
		ReadTimeout:       time.Second * 30,
		WriteTimeout:      time.Second * 10,
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
