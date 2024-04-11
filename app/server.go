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
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	directory := flag.String("directory", "", "user specified absolute path for the directory to serve")
	flag.Parse()
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go serve(conn, directory)
	}

}

func serve(conn net.Conn, directory *string) {
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("Error reading request:", err)
	}
	if req.Method == "GET" {
		handleGet(req, conn, directory)
	} else if req.Method == "POST" {
		handlePost(req, conn, directory)
	}
}

func handlePost(req *http.Request, conn net.Conn, directory *string) {
	uri := req.RequestURI
	if strings.HasPrefix(uri, "/files/") {
		splitUri := strings.SplitAfterN(uri, "/files/", 2)
		fileName := splitUri[1]
		absolutePath := *directory + "/" + fileName
		fileContent, _ := io.ReadAll(req.Body)
		err := os.WriteFile(absolutePath, fileContent, 0644)
		if err != nil {
			response := &http.Response{
				Status:     "404 Conflict",
				StatusCode: 409,
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
			response.Write(conn)
		} else {
			response := &http.Response{
				Status:     "201 OK",
				StatusCode: 201,
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
			response.Write(conn)
		}
	}
}

func handleGet(req *http.Request, conn net.Conn, directory *string) {
	uri := req.RequestURI
	fmt.Println("Request URI is:", uri)
	if uri == "/" {
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
	} else if strings.HasPrefix(uri, "/files/") {
		splitUri := strings.SplitAfterN(uri, "/files/", 2)
		fileName := splitUri[1]
		absolutePath := *directory + "/" + fileName
		fileContent, err := os.ReadFile(absolutePath)
		if err != nil {
			response := &http.Response{
				Status:     "404 Not Found",
				StatusCode: 404,
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
			response.Write(conn)
		} else {
			response := &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Type": []string{"application/octet-stream"},
				},
				Body:          io.NopCloser(bytes.NewReader(fileContent)),
				ContentLength: int64(len(fileContent)),
			}
			response.Write(conn)
		}
	} else if strings.HasPrefix(uri, "/user-agent") {
		userAgent := req.Header.Get("User-Agent")
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Content-Type": []string{"text/plain"},
			},
			Body:          io.NopCloser(strings.NewReader(userAgent)),
			ContentLength: int64(len(userAgent)),
		}
		response.Write(conn)
	} else if strings.HasPrefix(uri, "/echo/") {
		splitUri := strings.SplitAfterN(uri, "/echo/", 2)
		toEcho := splitUri[1]
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Content-Type": []string{"text/plain"},
			},
			Body:          io.NopCloser(strings.NewReader(toEcho)),
			ContentLength: int64(len(toEcho)),
		}
		response.Write(conn)
	} else {
		response := &http.Response{
			Status:     "404 Not Found",
			StatusCode: 404,
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
	}
}
