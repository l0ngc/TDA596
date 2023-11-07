package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

var listeningPort string
var maxProcesses int

func main() {
	if len(os.Args) < 2 {
		listeningPort = ":8080"
	} else {
		listeningPort = os.Args[1]
	}

	fmt.Printf("Listening on port%v\n", listeningPort)

	listener, err := net.Listen("tcp", listeningPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Connection initialized successfully")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue // handle error appropriately
		}

		fmt.Println("Connection established")
		go requestHandler(conn)
	}
}

func requestHandler(conn net.Conn) {
	// echo -e "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n" | nc localhost 8083
	// echo -e "POST /somepath HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 11\r\n\r\nhello=world" | nc localhost 8083
	defer conn.Close()
	fmt.Println("Handling request...")

	// init one buffer, and read info from the connection
	br := bufio.NewReaderSize(conn, 50*1024*1024) // 50MB buffer
	request, err := http.ReadRequest(br)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	// debug, now the basic connection is handled, then specify GET and POST handler for different request
	fmt.Println("Type of connection request: ", request.Method)
	switch request.Method {
	case "GET":
		getHandler(request)
	case "POST":
		postHandler(request)
	}

}

func getHandler(request *http.Request) {
	fmt.Println("Handler for get message")
}

func postHandler(request *http.Request) {
	fmt.Println("Handler for post message")
}
