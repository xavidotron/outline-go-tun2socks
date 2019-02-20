package tun2socks

import (
	"runtime/debug"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy"

	// "github.com/eycorsican/go-tun2socks/proxy/shadowsocks"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack

func init() {
	debug.SetGCPercent(20)
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			debug.FreeOSMemory()
		}
	}()
}

func InputPacket(data []byte) {
	lwipStack.Write(data)
}

func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int) {
	if packetFlow != nil {
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnectionHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
		core.RegisterUDPConnectionHandler(socks.NewUDPHandler(proxyHost, uint16(proxyPort), 30*time.Second, proxy.NewDNSCache()))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
	}
}

// func StartShadowsocks(packetFlow PacketFlow, proxyHost string, proxyPort int, proxyCipher, proxyPassword string) {
// 	if packetFlow != nil {
// 		lwipStack = core.NewLWIPStack()
// 		core.RegisterTCPConnectionHandler(shadowsocks.NewTCPHandler(core.ParseTCPAddr(proxyHost, uint16(proxyPort)).String(), proxyCipher, proxyPassword))
// 		core.RegisterUDPConnectionHandler(shadowsocks.NewUDPHandler(core.ParseUDPAddr(proxyHost, uint16(proxyPort)).String(), proxyCipher, proxyPassword, 30*time.Second))
// 		core.RegisterOutputFn(func(data []byte) (int, error) {
// 			packetFlow.WritePacket(data)
// 			return len(data), nil
// 		})
// 	}
// }
