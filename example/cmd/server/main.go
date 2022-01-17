package main

import (
	"flag"
	"log"
	socket "sockunx"
	"sockunx/example/handler"
)

var (
	socketPath = flag.String("socketPath", "/root/share/socks.sock", "specify socket path")
)

func main() {
	flag.Parse()

	log.Println("Starting server ...")

	server, e := socket.NewServer(*socketPath, 512)
	if e != nil {
		log.Fatal(e)
	}
	defer func() {
		log.Println("Shutting down...")
		server.Stop()
	}()

	// listening server
	log.Println("Server is started, waiting connection")
	server.RegisterHandler(handler.Index)
	e = server.Run()
	if e != nil {
		log.Println("Error while running server", e.Error())
	}
}
