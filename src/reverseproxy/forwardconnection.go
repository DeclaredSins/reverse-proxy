package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

func forwardtoclient(frontendconn net.Conn, backendconn net.Conn, rbuffer *bytes.Buffer, ch chan bool) {
	start := time.Now()
	tee := io.TeeReader(backendconn, rbuffer)
	n, err := io.Copy(frontendconn, tee)
	if err != nil {
		log.Println(n, "ERROR: ", err)
	}
	ch <- true
	elapsed := time.Since(start)
	log.Println("toclient", elapsed)
}

func forwardtoserver(frontendconn net.Conn, backendconn net.Conn, rbuffer *bytes.Buffer, ch chan bool) {
	start := time.Now()
	tee := io.TeeReader(frontendconn, rbuffer)
	n, err := io.Copy(backendconn, tee)
	if err != nil {
		log.Println(n, "ERROR: ", err)
	}
	ch <- true
	elapsed := time.Since(start)
	log.Println("toserver", elapsed)
}
