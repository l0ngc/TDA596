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
var stringfileTypes = []string{"html", "txt", "gif", "jpeg", "jpg", "css"}

func main() {
	// read port
	if len(os.Args) < 2 {
		listeningPort = ":8083"
	} else {
		listeningPort = os.Args[1]
	}

	fmt.Printf("Listening on port%v\n", listeningPort)
	// init tcp listening
	listener, err := net.Listen("tcp", listeningPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Connection initialized successfully")
	// init listening connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Port Listen Failed:", err)
			continue
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
	// distribute request to target handler
	switch request.Method {
	case "GET":
		getHandler(conn, request)
	case "POST":
		postHandler(conn, request)
	default:
		otherHandler(conn, request)
	}
}

func otherHandler(conn net.Conn, request *http.Request) {
	response := createResponse(http.StatusNotImplemented, "Not Implemented!")
	response.Write(conn)
}

// Design of Handler
// Inputs:
// 1. connection: response should be write back from the same coonection channel
// 2. Request pointer: for acquire the content of the request
// Outputs:
// 1. return the response code, represending whether success or not of this connect application
func getHandler(conn net.Conn, request *http.Request) {
	// echo -e "GET /resource/txt/monk.txt HTTP/1.1\r\nHost: localhost\r\n\r\n" | nc localhost 8083
	response := getResponseWrapper(request)
	response.Write(conn)
}

func getResponseWrapper(request *http.Request) *http.Response {
	// wrap up response body based on the Request
	url := request.URL.Path
	localPath, _ := os.Getwd()
	fileServerPath := localPath + url
	lastDotIndex := strings.LastIndex(url, ".") + 1

	contentType := url[lastDotIndex:]

	if debug == "true" {
		fmt.Println("The url PATH of request: ", url)
		fmt.Println("Current contentType is: ", contentType)
		fmt.Println("The current workdir of server: ", localPath)
		fmt.Println("Global path of target file: ", fileServerPath)
	}
	if isContentTypeSupported(contentType) {
		content, err := os.ReadFile(fileServerPath)
		if err != nil {
			response := createResponse(http.StatusNotFound, "File not Found!")
			return response
		}
		fileContent := string(content)
		response := createResponse(http.StatusOK, fileContent)
		return response
	} else {
		response := createResponse(http.StatusBadRequest, "Bad Request!")
		return response
	}
}

func postHandler(conn net.Conn, request *http.Request) {
	fmt.Println("Handler for post message")
	// echo -e "POST /download/test HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 11\r\n\r\nhello,world!" | nc localhost 8083

	response := postResponseWrapper(request)
	response.Write(conn)
}

func postResponseWrapper(request *http.Request) *http.Response {
	url := request.URL.Path
	localPath, _ := os.Getwd()
	fileSavePath := localPath + url
	fmt.Println("url: ", fileSavePath)
	// debug
	if debug == "true" {
		fmt.Println("The url PATH of request: ", url)
		fmt.Println("The current workdir of server: ", localPath)
		fmt.Println("Global path of target file: ", fileSavePath)
	}
	lastDotIndex := strings.LastIndex(url, ".") + 1
	contentType := url[lastDotIndex:]
	if isContentTypeSupported(contentType) {
		content, err := os.Create(fileSavePath)
		defer content.Close()
		if err != nil {
			fmt.Println("local file created failed")
			response := createResponse(http.StatusInternalServerError, "File not created successfull on server")
			return response
		}

		// Copy the response body to the local file.
		_, err = io.Copy(content, request.Body)

		if err != nil {
			fmt.Println("Error copying response to file:", err)
			response := createResponse(http.StatusInternalServerError, "File not created successfull on server")
			return response
		}
		response := createResponse(http.StatusCreated, "File created successfully")
		return response
	} else {
		response := createResponse(http.StatusBadRequest, "Bad Request!")
		return response
	}
	// create local file

}

func createResponse(statusCode int, message string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(message)),
	}
}

func isContentTypeSupported(contentType string) bool {
	for _, fileType := range stringfileTypes {
		if strings.Contains(contentType, fileType) {
			return true
		}
	}
	return false
}
