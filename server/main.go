package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)

	p, err := NewListener(c)
	if err != nil {
		log.Fatalln(err)
	}
	p.Listen()
}

