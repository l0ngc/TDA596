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
// func getHandler(conn net.Conn, request *http.Request) {
// 	content := "I want to fly"
// 	response := request.Response
// 	response.StatusCode = http.StatusOK
// 	response.ContentLength = 50
// 	response.Body = io.NopCloser(strings.NewReader(string(content)))

// }
func getHandler(conn net.Conn, request *http.Request) {

	// 创建一个新的 HTTP 响应
	// response := &http.Response{
	// 	Status:     http.StatusText(http.StatusOK),
	// 	StatusCode: http.StatusOK,
	// 	Body:       ioutil.NopCloser(strings.NewReader(content)),
	// }

	response := getResponseWrapper(request)
	// 将响应写回连接
	response.Write(conn)
}

func getResponseWrapper(request *http.Request) *http.Response {
	// wrap up one response body
	url := request.URL.Path
	localPath, _ := os.Getwd()
	fileServerPath := localPath + url

	content, err := os.ReadFile(fileServerPath)
	fileContent := string(content)

	// debug
	if debug == "true" {
		fmt.Println("The url PATH of request: ", url)
		fmt.Println("The current workdir of server: ", localPath)
		fmt.Println("Global path of target file: ", fileServerPath)
		fmt.Println("err: ", err)
	}

	if err != nil {
		// fault return
		response := &http.Response{
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("File not found")),
		}
		return response
	}

	// right return
	response := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(fileContent)),
	}

	return response
}

func postHandler(conn net.Conn, request *http.Request) {
	fmt.Println("Handler for post message")
}
