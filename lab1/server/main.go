package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

var listeningPort string
var maxProcesses int
var debug string // Global option for setting debug or not

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

// Design of Handler
// Inputs:
// 1. writer buffer: for writing back the response, message or data
// 2. Request pointer: for acquire the content of the request
// Outputs:
// 1. return the response code, represending whether success or not of this connect application
func getHandler(request *http.Request) {
	// echo -e "GET /resource/ebooks/monk.txt HTTP/1.1\r\nHost: localhost\r\n\r\n" | nc localhost 8083
	// TODO: Return the targeted request resource for users
	fmt.Println("Handler for get message")

	url := request.URL.Path
	local_path, _ := os.Getwd()
	global_path := local_path + url

	if debug == "true" {
		fmt.Println("The url PATH of request: ", url)
		fmt.Println("The current workdir of server: ", local_path)
		fmt.Println("Global path of target file: ", global_path)
	}

	file, _ := os.Open(global_path)
	defer file.Close()
	content, err := io.ReadAll(file)
	fileContent := string(content)
	fmt.Println("File content: ", fileContent)
	fmt.Println("Read error code: ", err)

	// http.ServeFile(request.Context().ResponseWriter, request, global_path)
	// now read the file on the global_path and return it by serveFile of http package,
}

func postHandler(request *http.Request) {
	fmt.Println("Handler for post message")
}
