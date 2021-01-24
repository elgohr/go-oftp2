package main

import (
	oftp2 "bifroest/oftp2/cmd"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Listener struct {
	c <-chan os.Signal
	listener *net.TCPListener
	connections map[string]struct{}
}

func NewListener(c <-chan os.Signal) (*Listener, error) {
	localAddress, err := net.ResolveTCPAddr("server", "0.0.0.0:3305")
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("server", localAddress)
	if err != nil {
		return nil, err
	}
	return &Listener{
		c:        c,
		listener: listener,
		connections: map[string]struct{}{},
	}, nil
}

func (p *Listener) Listen() {
	for {
		select {
		case <-p.c:
			log.Println("exiting...")
			return
		default:
			localConnection, err := p.listener.AcceptTCP()
			if err != nil {
				log.Println(err)
			}
			go p.handle(localConnection)
		}
	}
}

func (p *Listener) handle(connection *net.TCPConn) {
	reader := bufio.NewReader(connection)
	for {
		if _, exists := p.connections[connection.RemoteAddr().String()]; !exists {
			fmt.Printf("Serving %s\n", connection.RemoteAddr().String())
			if _, err := connection.Write(oftp2.StartSessionReadyMessage().StreamTransmissionBuffer()); err != nil {
				log.Println(err)
				return
			}
			//p.connections[connection.RemoteAddr().String()] = struct{}{}
		}

		c,_,  err := reader.ReadRune()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(c))
	}
}
