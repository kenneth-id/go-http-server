package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	directory := flag.String("directory", ".", "the directory to serve")
	flag.Parse()

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on port 4221")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, *directory)
	}
}

func handleConnection(conn net.Conn, directory string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	switch req.Method {
	case "GET":
		handleGet(conn, req, directory)
	case "POST":
		handlePost(conn, req, directory)
	default:
		sendResponse(conn, "405 Method Not Allowed", http.StatusMethodNotAllowed, nil, "")
	}
}

func handlePost(conn net.Conn, req *http.Request, directory string) {
	uri := req.RequestURI
	if strings.HasPrefix(uri, "/files/") {
		filePath := strings.TrimPrefix(uri, "/files/")
		absolutePath := directory + "/" + filePath
		fileContent, _ := io.ReadAll(req.Body)
		err := os.WriteFile(absolutePath, fileContent, 0644)
		if err != nil {
			sendResponse(conn, "409 Conflict", http.StatusConflict, nil, "")
		} else {
			sendResponse(conn, "201 Created", http.StatusCreated, nil, "")
		}
	}
}

func handleGet(conn net.Conn, req *http.Request, directory string) {
	uri := req.RequestURI
	switch {
	case uri == "/":
		sendResponse(conn, "200 OK", http.StatusOK, nil, "")
	case strings.HasPrefix(uri, "/files/"):
		filePath := strings.TrimPrefix(uri, "/files/")
		absolutePath := directory + "/" + filePath
		fileContent, err := os.ReadFile(absolutePath)
		if err != nil {
			sendResponse(conn, "404 Not Found", http.StatusNotFound, nil, "")
			return
		}
		sendResponse(conn, "200 OK", http.StatusOK, fileContent, "application/octet-stream")
	case strings.HasPrefix(uri, "/user-agent"):
		userAgent := req.Header.Get("User-Agent")
		sendResponse(conn, "200 OK", http.StatusOK, []byte(userAgent), "text/plain")
	case strings.HasPrefix(uri, "/echo/"):
		echoContent := strings.TrimPrefix(uri, "/echo/")
		sendResponse(conn, "200 OK", http.StatusOK, []byte(echoContent), "text/plain")
	default:
		sendResponse(conn, "404 Not Found", http.StatusNotFound, nil, "")
	}
}

func sendResponse(conn net.Conn, status string, statusCode int, body []byte, contentType string) {
	var bodyReader io.ReadCloser
	if body != nil {
		bodyReader = io.NopCloser(bytes.NewReader(body))
	}

	response := &http.Response{
		Status:        status,
		StatusCode:    statusCode,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          bodyReader,
		ContentLength: int64(len(body)),
	}

	if contentType != "" {
		response.Header.Set("Content-Type", contentType)
	}

	response.Write(conn)
}
