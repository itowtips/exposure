package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync"
)

var conf *config
var err error
var tunnelService string
var tunnelTcpAddr *net.TCPAddr
var tunnelListener *net.TCPListener
var tunnelConn *net.TCPConn
var mu sync.Mutex

func main() {
	conf, err = loadConfig()
	checkError(err)

	ListenTunnelService()
	ListenWebService()
}

func ListenTunnelService() {
	tunnelService = (*conf).TunnelService
	tunnelTcpAddr, err = net.ResolveTCPAddr("tcp4", tunnelService)
	checkError(err)
	tunnelListener, err = net.ListenTCP("tcp", tunnelTcpAddr)
	checkError(err)
	fmt.Printf("accepting tunnel connection ... %s\n", tunnelService)
	tunnelConn, err = tunnelListener.AcceptTCP()
	checkError(err)
}

func ListenWebService() {
	service := (*conf).FrontService
	fmt.Printf("Listen Web Service ... %s\n", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn *net.TCPConn) {
	defer func() {
		(*conn).Close()
	}()

	var reader *bufio.Reader
	var req *http.Request
	var res *http.Response

	// Read Request from FrontService
	fmt.Printf("(%p): read request from FrontService\n", conn)
	reader = bufio.NewReader(conn)
	req, err = http.ReadRequest(reader)
	checkError(err)

	mu.Lock()
	defer func() {
		mu.Unlock()
	}()

	// Write Request to HttpTunnel
	fmt.Printf("(%p): write request to HttpTunnel\n", conn)
	err = req.Write(tunnelConn)
	checkError(err)

	// Read Response from HttpTunnel
	fmt.Printf("(%p): read response from HttpTunnel\n", conn)
	reader = bufio.NewReader(tunnelConn)
	res, err = http.ReadResponse(reader, req)
	checkError(err)

	// Write Response to FrontService
	fmt.Printf("(%p): write response to FrontService\n", conn)
	err = res.Write(conn)
	checkError(err)
}
