package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/eycorsican/go-tun2socks/common/dns/cache"
	"github.com/eycorsican/go-tun2socks/core"

	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

const (
	mtu        = 1500
	bufferSize = 512 * 1024
)

// Usage:
// ./go-tun2socks-macos -proxyHost 127.0.0.1 -proxyPort 9999
//   -inboundSocketPath "~/Library/Containers/org.outline.macos.client.VpnExtension/Data/out_socket"
//   -outboundSocketPath "~/Library/Containers/org.outline.macos.client.VpnExtension/Data/in_socket"
func main() {
	proxyHost := flag.String("proxyHost", "", "proxy host")
	proxyPort := flag.Int("proxyPort", -1, "proxy port")
	inboundSocketPath := flag.String("inboundSocketPath", "", "inbound Unix socket path")
	outboundSocketPath := flag.String("outboundSocketPath", "", "outbound Unix socket path")
	flag.Parse()

	if *inboundSocketPath == "" || *outboundSocketPath == "" {
		fmt.Println("Must provide in/out Unix socket paths")
		os.Exit(1)
	}
	if *proxyHost == "" || *proxyPort == -1 {
		fmt.Println("Must provide a proxy host and port")
		os.Exit(1)
	}
	const connType = "unixgram"
	inAddr := net.UnixAddr{Name: *inboundSocketPath, Net: connType}
	outAddr := net.UnixAddr{Name: *outboundSocketPath, Net: connType}
	conn, err := net.DialUnix(connType, &inAddr, &outAddr)
	if err != nil {
		fmt.Printf("Failed to connect to socket: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	var lwipStack = core.NewLWIPStack()
	core.RegisterTCPConnHandler(socks.NewTCPHandler(*proxyHost, uint16(*proxyPort)))
	core.RegisterUDPConnHandler(socks.NewUDPHandler(*proxyHost, uint16(*proxyPort), 30*time.Second, cache.NewSimpleDnsCache()))
	core.RegisterOutputFn(func(data []byte) (int, error) {
		len, err := conn.Write(data)
		if err != nil {
			fmt.Printf("Failed to write packet %v\n", err.Error())
		}
		return len, err
	})
	conn.SetReadBuffer(bufferSize)
	conn.SetReadBuffer(bufferSize)
	var buf = make([]byte, mtu)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Failed to read packet %v\n", err.Error())
			continue
		}
		lwipStack.Write(buf)
	}
}
