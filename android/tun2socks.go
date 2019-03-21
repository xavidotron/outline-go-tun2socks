package tun2socks

import (
	"errors"
	"log"
	"os"

	"github.com/Jigsaw-Code/outline-go-tun2socks/common"
	"github.com/eycorsican/go-tun2socks/core"
)

const vpnMtu = 1500

var lwipStack core.LWIPStack
var tun *os.File
var config *common.ConnectionConfig
var isRunning = false

func StartSocks(fd int, proxyHost string, proxyPort int, isUDPSupported bool) error {
	lwipStack = core.NewLWIPStack()
	tun = os.NewFile(uintptr(fd), "")
	if tun == nil {
		return errors.New("Failed to open tun file descriptor")
	}
	config = &common.ConnectionConfig{
		Host: proxyHost, Port: uint16(proxyPort), IsUDPSupported: isUDPSupported}
	common.RegisterConnectionHandlers(config)
	core.RegisterOutputFn(func(data []byte) (int, error) {
		return tun.Write(data)
	})
	isRunning = true
	go processInputPackets()
	return nil
}

func StopSocks() {
	isRunning = false
	if tun != nil {
		tun.Close()
	}
	if lwipStack != nil {
		lwipStack.Close()
	}
}

func SetUDPSupport(isUDPSupported bool) error {
	if config.IsUDPSupported == isUDPSupported {
		return nil
	}
	config.IsUDPSupported = isUDPSupported
	if lwipStack != nil {
		lwipStack.Close() // Abort existing connections
	}
	lwipStack = core.NewLWIPStack()
	return common.RegisterConnectionHandlers(config)
}

func processInputPackets() {
	buffer := make([]byte, vpnMtu)
	for isRunning {
		len, err := tun.Read(buffer)
		if err != nil {
			log.Printf("Failed to read packet from TUN: %v", err)
			continue
		}
		if len == 0 {
			log.Println("Read EOF from TUN")
			continue
		}
		lwipStack.Write(buffer)
	}
}
