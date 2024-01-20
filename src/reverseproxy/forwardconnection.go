package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

func forwardtoclient(frontendconn net.Conn, backendconn net.Conn, rbuffer *bytes.Buffer, ch chan bool) {
	tee := io.TeeReader(backendconn, rbuffer)
	n, err := io.Copy(frontendconn, tee)
	if err != nil {
		log.Println(n, "ERROR: ", err)
	}
	ch <- true
}

func forwardtoserver(frontendconn net.Conn, backendconn net.Conn, rbuffer *bytes.Buffer, ch chan bool) {
	tee := io.TeeReader(frontendconn, rbuffer)
	n, err := io.Copy(backendconn, tee)
	if err != nil {
		log.Println(n, "ERROR: ", err)
	}
	ch <- true
}
