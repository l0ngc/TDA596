Your task is to build a web proxy capable of accepting HTTP requests, forwarding requests to remote (origin) servers, and returning response data to a client. The proxy will be implemented in Go and MUST handle concurrent requests by creating a Go routine for each new client request. You will only be responsible for implementing the GET method. All other request methods received by the proxy should elicit a "Not Implemented" (501) error (see RFC 1945Links to an external site. section 9.5 - Server Error). 


Steps:

1. Init port listening
2. check request works or not


echo -e "GET /resource/ebooks/monk.txt HTTP/1.1\r\nHost: localhost:8083\r\n\r\n" | nc localhost 8084