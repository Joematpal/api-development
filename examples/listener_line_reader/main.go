package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

const (
	contentLength = "Content-Length"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if err := ParseHTTP(conn); err != nil {
			log.Fatal(err)
		}
	}
}

func ParseHTTP(conn io.ReadWriteCloser) error {
	defer conn.Close()
	// conn.SetReadDeadline(time.Now().Add(time.Minute))

	// https: //developer.mozilla.org/en-US/docs/Glossary/Request_header
	reader := bufio.NewReader(io.LimitReader(conn, 1024))
	tp := textproto.NewReader(reader)

	var method string
	var path string
	var httpVersion string
	headers := map[string]string{}

	for {
		line, err := tp.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}

		if method == "" {
			split := strings.Split(line, " ")
			method = split[0]
			path = split[1]
			httpVersion = split[2]

			fmt.Println("method", method)
			fmt.Println("path", path)
			fmt.Println("httpVersion", httpVersion)
			continue
		}

		// Parse Headers
		val, ok := headers[contentLength]
		if !ok {
			split := strings.Split(line, ": ")
			headers[split[0]] = split[1]
			continue
		}

		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			continue
		}
		fmt.Println("content length: ", n)
		b := make([]byte, n)

		i, err := io.ReadAtLeast(tp.R, b, int(n))
		if err != nil {
			fmt.Println("err", err)
			continue
		}

		// Possible Headers
		// Connection: keep-alive
		// Content-Encoding: gzip
		// Content-Type: text/html; charset=utf-8
		// Date: Mon, 18 Jul 2016 16:06:00 GMT
		// Etag: "c561c68d0ba92bbeb8b0f612a9199f722e3a621a"
		// Keep-Alive: timeout=5, max=997
		// Last-Modified: Mon, 18 Jul 2016 02:36:04 GMT
		// Server: Apache
		// Set-Cookie: mykey=myvalue; expires=Mon, 17-Jul-2017 16:06:00 GMT; Max-Age=31449600; Path=/; secure
		// Transfer-Encoding: chunked
		// Vary: Cookie, Accept-Encoding
		// X-Backend-Server: developer2.webapp.scl3.mozilla.com
		// X-Cache-Info: not cacheable; meta data too large
		// X-kuma-revision: 1085259
		// x-frame-options: DENY

		// https://developer.mozilla.org/en-US/docs/Glossary/Response_header
		_, err = conn.Write([]byte(`HTTP/1.1 202 ACCEPTED
Access-Control-Allow-Origin: *
Connection: close
Conent-Type: text/plain
Content-Length: 0
`))
		if err != nil {
			fmt.Println("write error:", err)
		}
		// https://www.jmarshall.com/easy/http/#http1.1s4

		fmt.Println("bytes read:", i)
		fmt.Println("b", string(b))
		time.Sleep(time.Nanosecond * 10)
		return nil
	}
	return nil

}
