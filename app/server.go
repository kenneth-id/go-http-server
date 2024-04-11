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

// func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
// 	readBytes := []byte{}

// 	for {
// 		b, err := byteStream.ReadBytes('\n')
// 		if err != nil {
// 			return nil, err
// 		}

// 		readBytes = append(readBytes, b...)
// 		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
// 			break
// 		}
// 	}

// 	return readBytes[:len(readBytes)-2], nil
// }

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
		serve(conn)
	}

}

func serve(conn net.Conn) {
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("Error reading request:", err)
	}
	if req.Method == "GET" {
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
}
