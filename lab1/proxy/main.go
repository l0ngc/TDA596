package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

var listeningPort string
var maxProcesses int
var debug string // Global option for setting debug or not
// what
func main() {
	// read port
	if len(os.Args) < 2 {
		listeningPort = ":8083"
	} else {
		listeningPort = os.Args[1]
	}

	fmt.Printf("Start Proxy! Listening on port%v\n", listeningPort)
	// init tcp listening
	listener, err := net.Listen("tcp", listeningPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Connection initialized successfully")
	// because maximum connection is set on server side, no need to set it here
	for {
		client_conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Port Listen Failed:", err)
			continue
		}
		fmt.Println("Connection established")
		go proxyHandler(client_conn)
	}
}

// Design of Proxy Get Handler
// Inputs:
// 1. client connection: response should be write back from the same coonection channel
// 2. Request pointer: for acquire the content of the request
func proxyHandler(client_conn net.Conn) {
	defer client_conn.Close()
	fmt.Println("Proxy handling request...")

	// init one buffer, and read info from the connection
	br := bufio.NewReaderSize(client_conn, 50*1024*1024) // 50MB buffer
	request, err := http.ReadRequest(br)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	// debug, now the basic connection is handled, then specify GET and POST handler for different request
	fmt.Println("Type of connection request: ", request.Method)
	// init new connect with server
	// server_conn, _ := net.Dial("TCP", request.Host)
	// distribute request to target handler
	// only need to support GET request, others return StatusNotImplemented
	switch request.Method {
	case "GET":
		getHandler(client_conn, request)
	default:
		exceptHandler(client_conn, request)
	}

}

func exceptHandler(conn net.Conn, request *http.Request) {
	response := createResponse(http.StatusNotImplemented, "Request not supported")
	response.Write(conn)
}

// Design of Proxy Get Handler
// Inputs:
// 1. client connection: response should be write back from the same coonection channel
// 2. Request pointer: for acquire the content of the request
func getHandler(client_conn net.Conn, request *http.Request) {
	client_response := new(http.Response)
	// Forward request Host addr and Parameter to exact server
	server_response, err := http.Get("http://" + request.Host + request.URL.Path)
	if err != nil {
		// Return internal server error
		fmt.Println("Error reading server response:", err.Error())
		client_response = createResponse(http.StatusInternalServerError, "Request can not access Server")
		client_response.Write(client_conn)
		return
	}

	defer server_response.Body.Close()
	err = server_response.Write(client_conn)
	if err != nil {
		fmt.Println("Error writing client response:", err.Error())
		return
	}
}

func createResponse(statusCode int, message string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(message)),
	}
}
