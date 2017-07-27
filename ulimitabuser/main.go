package main

import (
	"flag"
	"log"
	"net"
)

var descriptors int

func init() {
	flag.IntVar(&descriptors, "n", 128, "how many descriptors to eat")
}

func main() {
	for i := 0; i < descriptors; i++ {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on %s", l.Addr())
	}
	log.Printf("Finished OK")
}
