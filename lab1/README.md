### Server build
cd ./server/
// enable debugging
go build -ldflags "-X main.debug=false" main.go
// without debugging
go build -ldflags "-X main.debug=true" main.go
### Start Server
./main 8080
### Proxy build
cd ./proxy/
go build main.go
### Start Server
./main 8081

### Listening 
When your server starts, the first thing that it will need to do is establish a socket connection that it can use to listen for incoming connections. Your server should listen on the port specified from the command line and wait for incoming client connections. Each new client request is accepted, and a new Go routine is spawned to handle the request. To avoid overwhelming your server, you should not create more than a reasonable number of child processes (for this assignment, use at most 10). In case an additional child process would break this limit, your server should wait until one of its ongoing child processes exits before forking a new one to handle the new request. 

Once a client has connected, the server should read data from the client and then check for a properly-formatted HTTP request. Your server should accept requests for files ending in html, txt, gif, jpeg, jpg, or css and transmit them to the client with a Content-Type of text/html, text/plain, image/gif, image/jpeg, image/jpeg, or text/css, respectively. If the client requests a file with any other extension, the web server must respond with a well-formed 400 "Bad Request" code. An invalid request from the client should be answered with an appropriate error code, i.e. "Bad Request" (400) or "Not Implemented" (501) for valid HTTP methods other than GET. If the requested file does not exist, your server should return a well-formed 404 "Not Found" code. Similarly, if headers are not properly formatted for parsing or any other error condition not listed before, your server should also generate a type-400 message.  For POST requests, please make sure that you store the files and make them accessible with a subsequent GET request.

### Parsing and Networking Libraries in Go 
For this assignment, you should use the package `net` for the networking, for example using `net.Listen("tcp", address)` to listen for incoming TCP connections. You can also use the package `net/http`, but ONLY for parsing and working with HTTP request objects, and not the networking part. You should not use e.g., `http.ListenAndServe` which trivializes the assignment (the same goes for `http.Listen`, and `http.Serve`). 

### Testing
There are no included tests for this assignment. However, to make sure your code works, you should set up some way of testing its functionalities yourself.


