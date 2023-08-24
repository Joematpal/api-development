package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
)

var (
	contentLengthRE = regexp.MustCompile(`Content-Length:\s`)
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
		defer conn.Close()
		fmt.Println("new conn")
		// conn.SetReadDeadline(time.Now().Add(time.Minute))
		reader := bufio.NewReader(io.LimitReader(conn, 1024))
		tp := textproto.NewReader(reader)

		var method string
		var path string
		var httpVersion string

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

			if contentLengthRE.MatchString(line) {
				lengthStr := contentLengthRE.ReplaceAllString(line, "")
				n, err := strconv.ParseInt(lengthStr, 10, 64)
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
				fmt.Println("bytes read:", i)
				fmt.Println("b", string(b))
			} else {
				// Parse Headers
				split := strings.Split(line, ": ")
				fmt.Println("Headers", split[0], split[1])
			}

		}
		fmt.Println("end")
	}
}
