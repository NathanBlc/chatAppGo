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

var isServerConnected = true

func handleMessage(c *net.TCPConn) {
	for {
		message := make([]byte, 2048)
		// message, _ := bufio.NewReader(c).ReadString('\n')
		_, err := c.Read(message)
		if err != nil && err == io.EOF {
			log.Printf("\nThe server stopped, stopping client too")
			isServerConnected = false
			break
		} else if err != nil {
			errorMessage := fmt.Errorf("An error occured while reading server response: %s", err.Error())
			fmt.Println(errorMessage)
			isServerConnected = false
			os.Exit(1)
		}
		fmt.Println("\nReceived from server ->: " + string(message))

		message = nil
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	// when leaving program close connection
	defer c.Close()
	// handle returning message
	go handleMessage(c.(*net.TCPConn))
	for isServerConnected {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Send a message to the server >> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
