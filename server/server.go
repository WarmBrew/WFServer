package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const ChunkSize = 1024 * 1024

var fileState = make(map[string]int64)

func main() {
	port := flag.String("port", "8080", "Specify the server port")
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("File server is listening on port %s...\n", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	infoBuf := make([]byte, 256)
	n, err := conn.Read(infoBuf)
	if err != nil {
		fmt.Println("Error reading file info:", err)
		return
	}
	info := strings.Split(string(infoBuf[:n]), "|")
	if len(info) < 2 {
		fmt.Println("Received incomplete file info")
		return
	}
	fileName := info[0]
	fileSize, _ := strconv.ParseInt(info[1], 10, 64)
	resume := info[2] == "true"

	offset := int64(0)
	if resume {
		offset = fileState[fileName]
		_, err = conn.Write([]byte(fmt.Sprintf("%d", offset)))
		if err != nil {
			fmt.Println("Error sending resume offset:", err)
			return
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buf := make([]byte, ChunkSize)
	var receivedBytes int64
	for offset < fileSize {
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading file chunk:", err)
			return
		}

		file.Seek(offset, 0)
		_, err = file.Write(buf[:n])
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		offset += int64(n)
		receivedBytes += int64(n)
		fileState[fileName] = offset
		fmt.Printf("Receiving %s: %d/%d bytes received\r", fileName, offset, fileSize)
	}

	fmt.Printf("\nFile received: %s, Size: %d bytes\n", fileName, receivedBytes)
	fmt.Println("File transfer completed successfully:", fileName)
}
