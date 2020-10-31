package main

import (
	"net"
	"strings"
)

func GetLocalIP() string {
	addresses, _ := net.InterfaceAddrs()
	for _, addr := range addresses {
		if !strings.HasPrefix(addr.String(), "127.") && !strings.Contains(addr.String(), ":") {
			return strings.TrimSuffix(addr.String(), "/24")
		}
	}
	return "Unknown"
}
