package main

import (
	"flag"
	"fmt"
	"net"
)

func startServer(portNum int) {
	fmt.Printf("Starting server on port %d ... \n", portNum)
	
	// Listen for a connection
	formatedPort := fmt.Sprintf(":%d", portNum)
    listener, error := net.Listen("tcp", formatedPort)
	
	if error != nil {
        fmt.Println(error)
	}
	
    manager := ConnectionManager{
        clients:    make(map[*Connection]bool),
        register:   make(chan *Connection),
	}
	
    go manager.Start()
	
	for {
        connection, _ := listener.Accept()
		
		if error != nil {
            fmt.Println(error)
		}
		
        client := &Connection{
			socket: connection,
			data: make(chan []byte)}
        manager.register <- client
		
		go manager.Receive(client)
        go manager.Send(client)
    }
}

// Main function that parses the cli port input
// Example: ./main --port 1234
// or you can run it without a port number which will default to port 8080
func main() {
	// Default port
	PORT := 8080

	// handle CLI input
	flagPort := flag.Int("port", PORT, "specify the server port, default is 8080")
	flag.Parse()
	PORT = *flagPort
	startServer(PORT)
}