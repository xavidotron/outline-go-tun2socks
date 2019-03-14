package tun2socks

import (
	"log"
	"os"
	"time"

	"github.com/eycorsican/go-tun2socks/common/dns/cache"
	"github.com/eycorsican/go-tun2socks/core"

	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

var lwipStack core.LWIPStack
var tun *os.File
var isRunning = false

func StartSocks(fd int, proxyHost string, proxyPort int) {
	lwipStack = core.NewLWIPStack()
	tun = os.NewFile(uintptr(fd), "")
	if tun == nil {
		log.Println("Failed to open tun file descriptor")
		return
	}
	core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
	core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, uint16(proxyPort), 30*time.Second, cache.NewSimpleDnsCache()))
	core.RegisterOutputFn(func(data []byte) (int, error) {
		return tun.Write(data)
	})
	isRunning = true
	go processInputPackets()
}

func StopSocks() {
	isRunning = false
	if lwipStack != nil {
		lwipStack.Close()
	}
}

func processInputPackets() {
	buffer := make([]byte, 1500)
	for isRunning {
		len, err := tun.Read(buffer)
		if err != nil {
			log.Println("Failed to read packet from TUN")
			continue
		}
		if len == 0 {
			log.Println("Read EOF from TUN")
			continue
		}
		lwipStack.Write(buffer)
	}
}
