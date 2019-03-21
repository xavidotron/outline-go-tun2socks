package tun2socks

import (
	"errors"
	"runtime/debug"
	"time"

	"github.com/Jigsaw-Code/outline-go-tun2socks/common"
	"github.com/eycorsican/go-tun2socks/core"
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack
var config *common.ConnectionConfig

func init() {
	// Apple VPN extensions have a memory limit of 15MB. Conserve memory by increasing garbage
	// collection frequency and returning memory to the OS every minute.
	debug.SetGCPercent(10)
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			debug.FreeOSMemory()
		}
	}()
}

func InputPacket(data []byte) {
	if lwipStack != nil {
		lwipStack.Write(data)
	}
}

func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int, isUDPSupported bool) error {
	if packetFlow == nil || proxyHost == "" || proxyPort <= 0 {
		return errors.New("Must provide a PacketFlow instance, valid proxy host and port")
	}
	config = &common.ConnectionConfig{
		Host: proxyHost, Port: uint16(proxyPort), IsUDPSupported: isUDPSupported}
	lwipStack = core.NewLWIPStack()
	common.RegisterConnectionHandlers(config)
	core.RegisterOutputFn(func(data []byte) (int, error) {
		packetFlow.WritePacket(data)
		return len(data), nil
	})
	return nil
}

func StopSocks() {
	if lwipStack != nil {
		lwipStack.Close()
	}
	lwipStack = nil
	config = nil
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
