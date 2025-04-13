package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChan := make(chan string)
	go func() {
		defer func() { fmt.Println("finished reading") }()
		defer f.Close()
		defer close(lineChan)

		fmt.Println("reading....")

		currentLine := ""
		for {
			chunk := make([]byte, 8)
			_, err := f.Read(chunk)
			if err == io.EOF {
				lineChan <- currentLine
				break
			}
			str := string(chunk)
			splittedStr := strings.Split(str, "\n")
			currentLine += splittedStr[0]

			if lenSplit := len(splittedStr); lenSplit > 1 {
				lineChan <- currentLine
				// aa\nbb\ncc -> [aa, bb, cc]. Till now aa has been sent now handle bb and cc
				for i := 1; i < lenSplit-1; i++ {
					lineChan <- splittedStr[i]
				}
				currentLine = splittedStr[lenSplit-1]
			}
		}

	}()
	return lineChan
}

func main() {
	tcp, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to setup tcp server")
		os.Exit(1)
	}
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Fprintln(os.Stdout, "error establishing connection")
		}
		fmt.Fprintln(os.Stdout, "Accepted connection from", conn.RemoteAddr())

		for line := range getLinesChannel(conn) {
			fmt.Fprintln(os.Stdout, line)
		}

		fmt.Fprintln(os.Stdout, "Connection to", conn.RemoteAddr(), "has been closed")
	}

}
