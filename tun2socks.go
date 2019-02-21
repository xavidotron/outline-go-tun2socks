package tun2socks

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy"

	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack

func init() {
	// Conserve memory in iOS by increasing garbage collection frequency and
	// returning memory to the OS every minute.
	if os.Getenv("GOOS") == "darwin" {
		debug.SetGCPercent(20)
		ticker := time.NewTicker(time.Minute * 1)
		go func() {
			for _ = range ticker.C {
				debug.FreeOSMemory()
			}
		}()
	}
}

func InputPacket(data []byte) {
	if lwipStack != nil {
		lwipStack.Write(data)
	}
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

func StopSocks() {
	if lwipStack != nil {
		lwipStack.Close()
	}
}
