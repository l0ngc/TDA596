package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var listeningPort string
var maxProcesses int

func main() {
	// TODO: init one listener, build one socket conenction

	if len(os.Args) < 2 {
		listeningPort = ":8080"
	} else {
		listeningPort = os.Args[1]
	}

	fmt.Printf("Port%v\n", listeningPort)

	listener, err := net.Listen("tcp", listeningPort)

	if err != nil {
		fmt.Println("Listen is err!: ", err)
	}

	defer listener.Close()

	fmt.Println("Conenction is inited successfully")

	// accept link
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Accept err!: ", err)
	}

	fmt.Println("Linking is inited successfully")
	go requestHandler(conn)
}

func requestHandler(conn net.Conn) {
	// defer conn.Close()
	fmt.Println("Connection established.")

	// Create a buffer to read data from the connection
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			break
		}
		// Output the received data
		fmt.Println("Received:", string(buffer[:n]))
		// You can also send data back to the client if needed
	}
}
