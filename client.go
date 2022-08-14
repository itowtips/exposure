package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

var conf *config
var err error
var tunnelService string
var tunnelTcpAddr *net.TCPAddr
var tunnelListener *net.TCPListener
var tunnelConn *net.TCPConn

func main() {
	conf, err = loadConfig()
	checkError(err)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run client.go common.go host:port\n")
		os.Exit(1)
	}
	tunnelService = os.Args[1]
	tunnelTcpAddr, err = net.ResolveTCPAddr("tcp4", tunnelService)
	checkError(err)
	tunnelConn, err = net.DialTCP("tcp", nil, tunnelTcpAddr)
	checkError(err)
	fmt.Println("established tunnel connection")

	var reader *bufio.Reader
	var req *http.Request
	var res *http.Response

	service := (*conf).BackendService
	var tcpAddr *net.TCPAddr
	var conn *net.TCPConn

	for {
		// Read Request from HttpTunnel
		fmt.Printf("read request from HttpTunnel\n")
		reader = bufio.NewReader(tunnelConn)
		req, err = http.ReadRequest(reader)
		checkError(err)

		// Write Request to BackendService
		fmt.Printf("write request to BackendService\n")
		tcpAddr, err = net.ResolveTCPAddr("tcp4", service)
		checkError(err)
		conn, err = net.DialTCP("tcp", nil, tcpAddr)
		checkError(err)
		err = req.Write(conn)
		checkError(err)

		// Read Response from BackendService
		fmt.Printf("read response from BackendService\n")
		reader = bufio.NewReader(conn)
		res, err = http.ReadResponse(reader, req)
		checkError(err)

		// Write Response to HttpTunnel
		fmt.Printf("write response to HttpTunnel\n")
		err = res.Write(tunnelConn)
		checkError(err)

		conn.Close()
	}
}
