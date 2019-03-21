package common

import (
	"errors"
	"time"

	"github.com/eycorsican/go-tun2socks/common/dns/cache"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/dnsfallback"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

// ConnectionConfig stores parameters of a connection to a proxy server.
type ConnectionConfig struct {
	Host           string
	Port           uint16
	IsUDPSupported bool
}

// RegisterConnectionHandlers registers UDP and TCP SOCKS connection handlers to the host and port
// specified in `config`. Registers a DNS/TCP fallback UDP handler when `config.IsUDPSupported` is
// false.
func RegisterConnectionHandlers(config *ConnectionConfig) error {
	if config == nil {
		return errors.New("Must provide a connection config")
	}
	var udpHandler core.UDPConnHandler
	if config.IsUDPSupported {
		udpHandler = socks.NewUDPHandler(
			config.Host, config.Port, 30*time.Second, cache.NewSimpleDnsCache(), nil)
	} else {
		udpHandler = dnsfallback.NewUDPHandler()
	}
	core.RegisterTCPConnHandler(socks.NewTCPHandler(config.Host, config.Port, nil))
	core.RegisterUDPConnHandler(udpHandler)
	return nil
}
