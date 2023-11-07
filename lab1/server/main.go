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
var debug string // Global option for setting debug or not

type response struct {
	status  string
	headers string
	body    string
}

func (r response) String() string {
	res := r.status + "\n" + r.headers + "\n" + r.body + "\n"
	return res
}

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
		getHandler(conn, request)
	case "POST":
		postHandler(conn, request)
	}

}

// Design of Handler
// Inputs:
// 1. connection: response should be write back from the same coonection channel
// 2. Request pointer: for acquire the content of the request
// Outputs:
// 1. return the response code, represending whether success or not of this connect application
func getHandler(conn net.Conn, request *http.Request) {
	// echo -e "GET /resource/ebooks/monk.txt HTTP/1.1\r\nHost: localhost\r\n\r\n" | nc localhost 8083
	// TODO: Return the targeted request resource for users
	fmt.Println("Handler for get message")

	tmp := getResponseWrapper(request)
	conn.Write([]byte(tmp))
}

func getResponseWrapper(request *http.Request) string {
	// wrap up one response body
	url := request.URL.Path
	local_path, _ := os.Getwd()
	file_server_path := local_path + url

	content, err := os.ReadFile(file_server_path)
	fileContent := string(content)
	if debug == "true" {
		fmt.Println("The url PATH of request: ", url)
		fmt.Println("The current workdir of server: ", local_path)
		fmt.Println("Global path of target file: ", file_server_path)
		fmt.Println("err: ", err)
	}

	// wrape one response body
	if err != nil {
		status := "HTTP/1.1 404 Not Found"
		header := "404:header"
		body := "404:body"
		resp := response{status, header, body}
		return resp.String()
	} else {
		status := "HTTP/1.1 404 Not Found"
		header := "404:header"
		body := fileContent
		resp := response{status, header, body}
		return resp.String()
	}
}

func postHandler(conn net.Conn, request *http.Request) {
	fmt.Println("Handler for post message")
}
