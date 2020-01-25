package base

import (
	"net"
	"regexp"
)

func CheckName(name string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name); !ok {
		return false
	}
	return true
}

func CheckIp(ip string) bool {
	return net.ParseIP(ip) != nil
}

func CheckIp4(ip string) bool {
	p := net.ParseIP(ip)
	if p != nil {
		return p.To4() != nil
	}

	return false
}

func ContainsAddress(address string, hostDetails []HostDetail) bool {
	for _, value := range hostDetails {
		if value.Address == address {
			return true
		}
	}
	return false
}
