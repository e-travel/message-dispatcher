package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/e-travel/message-dispatcher/dispatchers"
)

// TODO: what if this function panics?
func RequestHandler(conn io.Reader, buffer dispatchers.Dispatcher) {
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Print("Error reading from connection")
		return
	}
	// TODO: do some logging here if Put returns false
	buffer.Put(b)
}

func Serve(config *Config, handler func(io.Reader, dispatchers.Dispatcher), buffer dispatchers.Dispatcher, running chan<- bool) {

	// remove any existing socket file
	if config.socketType == "unix" {
		os.Remove(config.socketAddress)
	}
	listener, err := net.Listen(config.socketType, config.socketAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	running <- true

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(conn net.Conn, handler func(io.Reader, dispatchers.Dispatcher)) {
			handler(conn, buffer)
			conn.Close()
		}(conn, handler)
	}
}
