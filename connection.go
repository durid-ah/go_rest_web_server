package main

import (
	"net"
)

// Connection : Handles a single connection
type Connection struct {
    socket net.Conn 	// the client tcp connection
    data   chan []byte 	// response to be sent back
}
