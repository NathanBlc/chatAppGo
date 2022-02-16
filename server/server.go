package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// Type declaration
type Message struct {
	sender string
	msg    string
}

var count = 0

var clients = make(map[*net.TCPConn]bool)
var broadcaster = make(chan Message)

// handleConnection will handle all connection to the server
func handleConnection(c net.Conn) {
	// close connection when leaving handleConnection function
	defer c.Close()
	fmt.Print(".")
	isConnected := true
	for isConnected {
		// Handle connection close
		isConnected = connIsClosed(c.(*net.TCPConn))
	}
	// remove client from map to reach
	delete(clients, c.(*net.TCPConn))
}

// checkErr will check if there is an error and log it
func checkErr(err error) {
	if err != nil {
		log.Panicf("Error: server.go %s\n", err.Error())
	}
}

// handle message from channel and send it to all clients in clients
func handleMessage() {
	for {
		// Retrieve message from channel
		message := <-broadcaster
		for client := range clients {
			if client.RemoteAddr().String() != message.sender {
				// send message to all client connected
				client.Write([]byte(message.msg))
			}
		}
	}
}

// connIsClosed is a function which will detect when a connection is closed
func connIsClosed(c *net.TCPConn) (isConnected bool) {
	netData, err := bufio.NewReader(c).ReadString('\n')
	message := strings.TrimSpace(string(netData))
	if err == io.EOF || message == "STOP" {
		log.Printf("Client disconnect: %s \n", c.RemoteAddr())
		isConnected = false
		count--
	} else if message != "" {
		// print message received
		log.Printf("Message received from %s: %s", c.RemoteAddr().String(), message)
		// create Message to send to all other clients
		msgToSend := Message{sender: c.RemoteAddr().String(), msg: message}
		// send message through the channel
		broadcaster <- msgToSend
		isConnected = true
	}
	return
}

// logClientJoined is a function which will log when a client joined the server
func logClientJoined(conn *net.TCPConn) {
	log.Printf("server.go: Client joined from %s \n", conn.RemoteAddr())
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()
	// handle messages
	go handleMessage()

	for {
		c, err := l.Accept()
		checkErr(err)
		// Log client joined
		logClientJoined(c.(*net.TCPConn))
		// Add client to map
		clients[c.(*net.TCPConn)] = true
		// handle connection
		go handleConnection(c)
		count++
	}

}
