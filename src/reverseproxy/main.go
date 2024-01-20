package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Accept incoming connections
		frontendconn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle backend connection in a goroutine
		go handlebackend(frontendconn)
	}
}

func handlebackend(frontendconn net.Conn) {
	defer frontendconn.Close()

	log.Println("Connection from " + frontendconn.RemoteAddr().String())

	// Read and process data from the backend
	backendconn, err := net.DialTimeout("tcp", "192.168.70.2:31826", 3*time.Second)
	defer backendconn.Close()

	if err != nil {
		log.Fatalf("Unable to set backendConn deadline %v", err)
	}

	log.Print("frontendConnected")
	if err != nil {
		log.Println("Error:", err)
	}
	requestBuf := new(bytes.Buffer)
	responseBuf := new(bytes.Buffer)
	ch := make(chan bool)

	go forwardtoclient(frontendconn, backendconn, requestBuf, ch)
	// forward data from server to backend
	go forwardtoserver(frontendconn, backendconn, responseBuf, ch)

	<-ch
	<-ch
}
