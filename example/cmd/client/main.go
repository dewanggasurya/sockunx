package main

import (
	"flag"
	"fmt"
	"log"
	socket "sockunx"
)

var (
	socketPath = flag.String("socketPath", "/root/share/socks.sock", "specify socket path")
)

func main() {
	flag.Parse()

	client, e := socket.NewClient(*socketPath)
	if e != nil {
		log.Fatal(e)
	}

	for i := 0; i < 10; i++ {
		response, e := client.Send(`{"id":"one","from":0,"to":15,"fizz":"zzif","buzz":"zzub"}\n`)
		if e != nil {
			log.Println("error : ", e.Error())
			continue
		}
		fmt.Println(">>>", response)
	}
}
