package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	toml "github.com/pelletier/go-toml/v2"
)

type config struct {
	Destination_ipaddr string
	Port               string
}

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
			panic(err)
		}

		// Handle backend connection in a goroutine
		go handlebackend(frontendconn)
	}
}

func handlebackend(frontendconn net.Conn) {
	defer frontendconn.Close()

	log.Println("Connection from " + frontendconn.RemoteAddr().String())
	ipaddr := readconfig()
	log.Println(ipaddr)

	backendconn, err := net.DialTimeout("tcp", ipaddr, 3*time.Second)
	defer backendconn.Close()

	if err != nil {
		panic(err)
	}

	requestBuf := new(bytes.Buffer)
	responseBuf := new(bytes.Buffer)
	ch := make(chan bool)

	go forwardtoclient(frontendconn, backendconn, requestBuf, ch)
	go forwardtoserver(frontendconn, backendconn, responseBuf, ch)

	<-ch
	<-ch
}

func readconfig() string {
	path, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	currentpath := filepath.Dir(path)
	doc, err := os.ReadFile(currentpath + "\\config.toml")

	var cfg config
	errs := toml.Unmarshal(doc, &cfg)
	if errs != nil {
		panic(errs)
	}

	// lookup for dynamic ip
	ips, err := net.LookupIP(cfg.Destination_ipaddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
	}

	var resolvedip string
	for _, ip := range ips {
		resolvedip = ip.String()
	}

	log.Println(resolvedip)
	return resolvedip + ":" + cfg.Port
}
