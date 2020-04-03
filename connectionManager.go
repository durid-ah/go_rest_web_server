package main

import(
	"fmt"
)

// ConnectionManager : Handles incoming connections and their respective clients
type ConnectionManager struct {
    clients    map[*Connection]bool
    register   chan *Connection
}

// Start : Registers the clients into a map
func (manager *ConnectionManager) Start() {
    for {
        select {
        case connection := <-manager.register:
            manager.clients[connection] = true
            fmt.Println("Added new connection!")
        }
    }
}

// Send : A goroutine that handles the response package after the system gets it ready
func (manager *ConnectionManager) Send(client *Connection) {
    defer client.socket.Close()
    for {
        select {
        case message, ok := <-client.data:
            if !ok {
                return
            }
            client.socket.Write(message)
        }
    }
}

// Receive : Handles incoming requests from the client 
func (manager *ConnectionManager) Receive(client *Connection) {
    for {
        message := make([]byte, 4096)
        length, err := client.socket.Read(message)
		
		if err != nil {
            client.socket.Close()
            break
		}
		
        if length > 0 {
			response := ParseRequest(string(message), length) 
			client.data <- []byte(response)
		}

		close(client.data)
		delete(manager.clients, client)
    }
}
